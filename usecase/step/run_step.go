package step

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"

	"flowforge-api/domain/types"
	"flowforge-api/infrastructure/messaging/rabbitmq"
	"flowforge-api/infrastructure/runner"
	"flowforge-api/presenter"
)

type RunStepUseCase struct {
	httpClient *http.Client
}

func NewRunStepUseCase(httpClient *http.Client) *RunStepUseCase {
	return &RunStepUseCase{
		httpClient: httpClient,
	}
}

func (u *RunStepUseCase) Execute(ctx context.Context, stepRunEvent *rabbitmq.StepRunEvent) (runner.RunnerResponse, error) {
	queueTime := time.Duration(0)
	if stepRunEvent.QueuedAt != nil {
		queueTime = time.Since(*stepRunEvent.QueuedAt)
		if queueTime < 0 {
			queueTime = 0
		}
	}

	step := stepRunEvent.Step
	endpoint := stepRunEvent.Endpoint
	stepStartTime := time.Now()

	config := u.resolveConfig(step, endpoint)

	maxAttempts := 1
	if config.RetryOnFailure && config.RetryCount > 0 {
		maxAttempts = config.RetryCount + 1
	}

	var lastResponse runner.RunnerResponse
	aggregatedInsights := runner.RunnerInsights{
		TotalAttempts: maxAttempts,
		QueueTime:     queueTime,
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		attemptInsights := runner.RunnerInsights{
			AttemptNumber: attempt,
			TotalAttempts: maxAttempts,
			QueueTime:     queueTime,
		}

		response, err := u.executeRequest(config, &attemptInsights)
		lastResponse = response
		aggregatedInsights = u.aggregateInsights(aggregatedInsights, attemptInsights)

		if err == nil {
			aggregatedInsights.Duration = time.Since(stepStartTime)
			lastResponse.Insights = aggregatedInsights
			return lastResponse, nil
		}

		if attempt < maxAttempts && config.RetryDelay > 0 {
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
		}
	}

	aggregatedInsights.Duration = time.Since(stepStartTime)
	lastResponse.Insights = aggregatedInsights
	return lastResponse, fmt.Errorf("execution failed after %d attempts", maxAttempts)
}

func (u *RunStepUseCase) resolveConfig(step presenter.StepDetailResponse, endpoint presenter.EndpointDetailResponse) runner.ExecutionConfig {
	config := runner.ExecutionConfig{
		URL:    endpoint.BaseURI + endpoint.Path,
		Method: endpoint.Method,
	}

	config.Timeout = endpoint.Timeout
	if step.Timeout > 0 {
		config.Timeout = step.Timeout
	}

	if step.RetryOnFailure {
		config.RetryOnFailure = true
		config.RetryCount = step.RetryCount
		config.RetryDelay = step.RetryDelay
	} else if endpoint.RetryOnFailure {
		config.RetryOnFailure = true
		config.RetryCount = endpoint.RetryCount
		config.RetryDelay = endpoint.RetryDelay
	}

	config.URL = u.buildURL(config.URL, step.Query, endpoint.Query)

	config.Headers = u.mergeHeaders(endpoint.Header, step.Header)

	if u.hasValidBody(step.Body) {
		config.Body = step.Body
	} else if u.hasValidBody(endpoint.Body) {
		config.Body = endpoint.Body
	}

	return config
}

func (u *RunStepUseCase) buildURL(baseURL string, stepQuery, endpointQuery types.Query) string {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	queryParams := parsedURL.Query()

	for _, param := range endpointQuery {
		if param.Key != "" {
			queryParams.Add(param.Key, param.Value)
		}
	}

	for _, param := range stepQuery {
		if param.Key != "" {
			queryParams.Set(param.Key, param.Value)
		}
	}

	parsedURL.RawQuery = queryParams.Encode()
	return parsedURL.String()
}

func (u *RunStepUseCase) mergeHeaders(endpointHeaders, stepHeaders types.Header) http.Header {
	headers := make(http.Header)

	for _, h := range endpointHeaders {
		if h.Key != "" {
			headers.Set(h.Key, h.Value)
		}
	}

	for _, h := range stepHeaders {
		if h.Key != "" {
			headers.Set(h.Key, h.Value)
		}
	}

	return headers
}

func (u *RunStepUseCase) hasValidBody(body types.Body) bool {
	if len(body) == 0 {
		return false
	}

	bodyStr := string(body)
	return bodyStr != "null" && bodyStr != "[]" && bodyStr != "{}" && bodyStr != ""
}

func (u *RunStepUseCase) supportsBody(method string) bool {
	return method != "GET" && method != "HEAD" && method != "DELETE"
}

func (u *RunStepUseCase) executeRequest(config runner.ExecutionConfig, insights *runner.RunnerInsights) (runner.RunnerResponse, error) {

	var bodyReader io.Reader
	if len(config.Body) > 0 && u.supportsBody(config.Method) {
		bodyReader = bytes.NewReader(config.Body)
	}

	req, err := http.NewRequest(config.Method, config.URL, bodyReader)
	if err != nil {
		insights.ErrorType = "failed to create request"
		insights.ErrorMessage = err.Error()
		return runner.RunnerResponse{Insights: *insights}, fmt.Errorf("%s: %w", insights.ErrorType, err)
	}

	req.Header = config.Headers
	insights.RequestSize = int64(len(config.Body))

	var requestStartTime time.Time
	var dnsStartTime time.Time
	var connectStartTime time.Time
	var tlsStartTime time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			if dnsStartTime.IsZero() {
				dnsStartTime = time.Now()
			}
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			if !dnsStartTime.IsZero() {
				insights.DNSLookupDuration = time.Since(dnsStartTime)
			}
		},
		ConnectStart: func(_, _ string) {
			if connectStartTime.IsZero() {
				connectStartTime = time.Now()
			}
		},
		ConnectDone: func(_, _ string, _ error) {
			if !connectStartTime.IsZero() {
				insights.TCPConnectionTime = time.Since(connectStartTime)
			}
		},
		TLSHandshakeStart: func() {
			if tlsStartTime.IsZero() {
				tlsStartTime = time.Now()
			}
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			if !tlsStartTime.IsZero() {
				insights.TLSHandshakeTime = time.Since(tlsStartTime)
			}
		},
		GotFirstResponseByte: func() {
			if !requestStartTime.IsZero() {
				insights.TTFB = time.Since(requestStartTime)
			}
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	insights.StartTime = time.Now()
	requestStartTime = insights.StartTime
	resp, err := client.Do(req)
	insights.EndTime = time.Now()
	insights.Duration = insights.EndTime.Sub(insights.StartTime)

	if err != nil {
		insights.ErrorType = fmt.Sprintf("HTTP request failed (%s %s)", config.Method, config.URL)
		insights.ErrorMessage = err.Error()
		return runner.RunnerResponse{Insights: *insights}, fmt.Errorf("%s: %w", insights.ErrorType, err)
	}
	defer resp.Body.Close()

	insights.StatusCode = resp.StatusCode

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		insights.ErrorType = "failed to read response body"
		insights.ErrorMessage = err.Error()
		return runner.RunnerResponse{Insights: *insights}, fmt.Errorf("%s: %w", insights.ErrorType, err)
	}

	insights.ResponseSize = int64(len(respBody))

	if resp.StatusCode >= 400 {
		insights.ErrorType = fmt.Sprintf("HTTP %d response error", resp.StatusCode)
		insights.ErrorMessage = string(respBody)
		return runner.RunnerResponse{
			Response: string(respBody),
			Insights: *insights,
		}, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return runner.RunnerResponse{
		Response: string(respBody),
		Insights: *insights,
	}, nil
}

func (u *RunStepUseCase) aggregateInsights(total runner.RunnerInsights, attempt runner.RunnerInsights) runner.RunnerInsights {
	if total.StartTime.IsZero() || (!attempt.StartTime.IsZero() && attempt.StartTime.Before(total.StartTime)) {
		total.StartTime = attempt.StartTime
	}
	if attempt.EndTime.After(total.EndTime) {
		total.EndTime = attempt.EndTime
	}

	total.AttemptNumber = attempt.AttemptNumber
	total.StatusCode = attempt.StatusCode
	total.RequestSize += attempt.RequestSize
	total.ResponseSize += attempt.ResponseSize
	total.DNSLookupDuration += attempt.DNSLookupDuration
	total.TCPConnectionTime += attempt.TCPConnectionTime
	total.TLSHandshakeTime += attempt.TLSHandshakeTime
	total.TTFB += attempt.TTFB

	if attempt.ErrorType != "" {
		total.ErrorType = attempt.ErrorType
	}
	if attempt.ErrorMessage != "" {
		total.ErrorMessage = attempt.ErrorMessage
	}

	return total
}
