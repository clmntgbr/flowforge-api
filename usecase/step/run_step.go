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
	fmt.Printf("[run_step] Execute START step_run_id=%s workflow_run_id=%s workflow_id=%s\n",
		stepRunEvent.StepRunID, stepRunEvent.WorkflowRunID, stepRunEvent.WorkflowID)

	queueTime := time.Duration(0)
	if stepRunEvent.QueuedAt != nil {
		queueTime = time.Since(*stepRunEvent.QueuedAt)
		if queueTime < 0 {
			queueTime = 0
		}
		fmt.Printf("[run_step] queue_time=%s queued_at=%s\n", queueTime, stepRunEvent.QueuedAt.Format(time.RFC3339))
	} else {
		fmt.Println("[run_step] no QueuedAt on event")
	}

	step := stepRunEvent.Step
	endpoint := stepRunEvent.Endpoint
	stepStartTime := time.Now()

	config := ResolveConfig(step, endpoint)
	fmt.Printf("[run_step] config resolved method=%s url=%s timeout=%ds retry_on_failure=%v retry_count=%d retry_delay=%ds\n",
		config.Method, config.URL, config.Timeout, config.RetryOnFailure, config.RetryCount, config.RetryDelay)

	maxAttempts := 1
	if config.RetryOnFailure && config.RetryCount > 0 {
		maxAttempts = config.RetryCount + 1
	}
	fmt.Printf("[run_step] max_attempts=%d\n", maxAttempts)

	var lastResponse runner.RunnerResponse
	aggregatedInsights := runner.RunnerInsights{
		TotalAttempts: maxAttempts,
		QueueTime:     queueTime,
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("[run_step] attempt %d/%d START\n", attempt, maxAttempts)

		attemptInsights := runner.RunnerInsights{
			AttemptNumber: attempt,
			TotalAttempts: maxAttempts,
			QueueTime:     queueTime,
		}

		response, err := ExecuteRequest(config, &attemptInsights)
		lastResponse = response
		aggregatedInsights = AggregateInsights(aggregatedInsights, attemptInsights)

		if err == nil {
			fmt.Printf("[run_step] attempt %d/%d SUCCESS status=%d duration=%s response_len=%d\n",
				attempt, maxAttempts, attemptInsights.StatusCode, attemptInsights.Duration, len(response.Response))
			aggregatedInsights.Duration = time.Since(stepStartTime)
			lastResponse.Insights = aggregatedInsights
			fmt.Printf("[run_step] Execute DONE success total_duration=%s\n", aggregatedInsights.Duration)
			return lastResponse, nil
		}

		fmt.Printf("[run_step] attempt %d/%d FAILED error_type=%q error_msg=%q status=%d\n",
			attempt, maxAttempts, attemptInsights.ErrorType, attemptInsights.ErrorMessage, attemptInsights.StatusCode)

		if attempt < maxAttempts && config.RetryDelay > 0 {
			fmt.Printf("[run_step] sleeping %ds before retry\n", config.RetryDelay)
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
			fmt.Printf("[run_step] retry sleep done\n")
		}
	}

	aggregatedInsights.Duration = time.Since(stepStartTime)
	lastResponse.Insights = aggregatedInsights
	fmt.Printf("[run_step] Execute DONE failed after %d attempts total_duration=%s last_error=%q\n",
		maxAttempts, aggregatedInsights.Duration, aggregatedInsights.ErrorMessage)
	return lastResponse, fmt.Errorf("execution failed after %d attempts", maxAttempts)
}

func ResolveConfig(step presenter.StepDetailResponse, endpoint presenter.EndpointDetailResponse) runner.ExecutionConfig {
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

	config.URL = BuildURL(config.URL, step.Query, endpoint.Query)

	config.Headers = MergeHeaders(endpoint.Header, step.Header)

	if HasValidBody(step.Body) {
		config.Body = step.Body
	} else if HasValidBody(endpoint.Body) {
		config.Body = endpoint.Body
	}

	return config
}

func BuildURL(baseURL string, stepQuery, endpointQuery types.Query) string {
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

func MergeHeaders(endpointHeaders, stepHeaders types.Header) http.Header {
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

func HasValidBody(body types.Body) bool {
	if len(body) == 0 {
		return false
	}

	bodyStr := string(body)
	return bodyStr != "null" && bodyStr != "[]" && bodyStr != "{}" && bodyStr != ""
}

func SupportsBody(method string) bool {
	return method != "GET" && method != "HEAD" && method != "DELETE"
}

func ExecuteRequest(config runner.ExecutionConfig, insights *runner.RunnerInsights) (runner.RunnerResponse, error) {
	fmt.Printf("[run_step] ExecuteRequest START attempt=%d %s %s body_len=%d\n",
		insights.AttemptNumber, config.Method, config.URL, len(config.Body))

	var bodyReader io.Reader
	if len(config.Body) > 0 && SupportsBody(config.Method) {
		bodyReader = bytes.NewReader(config.Body)
	}

	req, err := http.NewRequest(config.Method, config.URL, bodyReader)
	if err != nil {
		fmt.Printf("[run_step] ExecuteRequest failed to create request: %v\n", err)
		insights.ErrorType = "failed to create request"
		insights.ErrorMessage = err.Error()
		return runner.RunnerResponse{Insights: *insights}, fmt.Errorf("%s: %w", insights.ErrorType, err)
	}

	req.Header = config.Headers
	insights.RequestSize = int64(len(config.Body))
	fmt.Printf("[run_step] ExecuteRequest sending HTTP request timeout=%ds\n", config.Timeout)

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
	fmt.Printf("[run_step] ExecuteRequest client.Do in progress...\n")
	resp, err := client.Do(req)
	insights.EndTime = time.Now()
	insights.Duration = insights.EndTime.Sub(insights.StartTime)
	fmt.Printf("[run_step] ExecuteRequest client.Do returned duration=%s err=%v\n", insights.Duration, err)

	if err != nil {
		fmt.Printf("[run_step] ExecuteRequest HTTP error: %v\n", err)
		insights.ErrorType = fmt.Sprintf("HTTP request failed (%s %s)", config.Method, config.URL)
		insights.ErrorMessage = err.Error()
		return runner.RunnerResponse{Insights: *insights}, fmt.Errorf("%s: %w", insights.ErrorType, err)
	}
	defer resp.Body.Close()

	insights.StatusCode = resp.StatusCode
	fmt.Printf("[run_step] ExecuteRequest status=%d reading body...\n", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[run_step] ExecuteRequest failed to read body: %v\n", err)
		insights.ErrorType = "failed to read response body"
		insights.ErrorMessage = err.Error()
		return runner.RunnerResponse{Insights: *insights}, fmt.Errorf("%s: %w", insights.ErrorType, err)
	}

	insights.ResponseSize = int64(len(respBody))
	fmt.Printf("[run_step] ExecuteRequest body read response_size=%d\n", insights.ResponseSize)

	if resp.StatusCode >= 400 {
		fmt.Printf("[run_step] ExecuteRequest HTTP error status=%d body_preview=%q\n", resp.StatusCode, truncateForLog(string(respBody), 200))
		insights.ErrorType = fmt.Sprintf("HTTP %d response error", resp.StatusCode)
		insights.ErrorMessage = string(respBody)
		return runner.RunnerResponse{
			Response: string(respBody),
			Insights: *insights,
		}, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	fmt.Printf("[run_step] ExecuteRequest DONE success status=%d ttfb=%s\n", resp.StatusCode, insights.TTFB)
	return runner.RunnerResponse{
		Response: string(respBody),
		Insights: *insights,
	}, nil
}

func truncateForLog(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func AggregateInsights(total runner.RunnerInsights, attempt runner.RunnerInsights) runner.RunnerInsights {
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
