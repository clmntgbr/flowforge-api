package mercure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Publisher struct {
	hubURL    string
	jwtSecret []byte
	client    *http.Client
}

func NewPublisher(hubURL string, jwtSecret string) *Publisher {
	return &Publisher{
		hubURL:    hubURL,
		jwtSecret: []byte(jwtSecret),
		client: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) > 0 && via[0].URL.Scheme != req.URL.Scheme {
					return fmt.Errorf("refusing redirect from %s to %s", via[0].URL, req.URL)
				}
				return nil
			},
		},
	}
}

func (p *Publisher) Publish(topic string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal mercure data: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"mercure": map[string]any{
			"publish": []string{topic},
		},
		"exp": time.Now().Add(time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(p.jwtSecret)
	if err != nil {
		return fmt.Errorf("sign mercure jwt: %w", err)
	}

	form := url.Values{}
	form.Set("topic", topic)
	form.Set("data", string(payload))

	req, err := http.NewRequest(http.MethodPost, p.hubURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("create mercure request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+tokenString)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("publish to mercure: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("mercure hub returned status %d", resp.StatusCode)
	}

	return nil
}
