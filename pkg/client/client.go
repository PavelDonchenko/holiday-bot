package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

//go:generate mockery --name=Client --output=mock --case=underscore
type Client interface {
	SendRequest(req *http.Request) (*http.Response, error)
	BuildURL(resource string, filters []FilterOptions) (string, error)
}
type BaseClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (c *BaseClient) SendRequest(req *http.Request) (*http.Response, error) {
	if c.HTTPClient == nil {
		return nil, errors.New("no http client")
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request. error: %w", err)
	}

	return response, nil
}

func (c *BaseClient) BuildURL(resource string, filters []FilterOptions) (string, error) {
	var resultURL string
	parsedURL, err := url.ParseRequestURI(c.BaseURL)
	if err != nil {
		return resultURL, fmt.Errorf("failed to parse base URL. error: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, resource)

	if len(filters) > 0 {
		q := parsedURL.Query()
		for _, fo := range filters {
			q.Set(fo.Field, fo.ToStringWF())
		}
		parsedURL.RawQuery = q.Encode()
	}

	return parsedURL.String(), nil
}

func (c *BaseClient) Close() error {
	c.HTTPClient = nil
	return nil
}
