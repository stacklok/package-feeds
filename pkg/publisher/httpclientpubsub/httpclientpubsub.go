package httpclientpubsub

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const PublisherType = "http-client"

type Config struct {
	URL string `mapstructure:"url"`
}

type HTTPClientPubSub struct {
	url string
}

func New(ctx context.Context, url string) (*HTTPClientPubSub, error) {
	pub := &HTTPClientPubSub{url: url}
	return pub, nil
}

func (pub *HTTPClientPubSub) Name() string {
	return PublisherType
}

func FromConfig(ctx context.Context, config Config) (*HTTPClientPubSub, error) {
	return New(ctx, config.URL)
}

func (pub *HTTPClientPubSub) Send(ctx context.Context, body []byte) error {
	log.Info("Sending event to HTTP client publisher")
	// Print the url to the log so that we can see where the event is being sent.
	req, err := http.NewRequest("POST", pub.url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	return nil
}
