package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type Error struct {
	Method     string
	Path       string
	StatusCode int
	Code       models.ErrorCode
	Message    string
}

func (e Error) Error() string {
	if e.Message != "" {
		if e.Code != "" {
			return fmt.Sprintf("%s %s: http %d: %s: %s", e.Method, e.Path, e.StatusCode, e.Code, e.Message)
		}
		return fmt.Sprintf("%s %s: http %d: %s", e.Method, e.Path, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("%s %s: http %d", e.Method, e.Path, e.StatusCode)
}

func New(baseURL string) (*Client, error) {
	baseURL, err := normalizeBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *Client) doJSON(ctx context.Context, method string, path string, out any) error {
	return c.doJSONBody(ctx, method, path, nil, out)
}

func (c *Client) doJSONBody(ctx context.Context, method string, path string, in any, out any) error {
	var body io.Reader
	if in != nil {
		payload, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("encode %s %s: %w", method, path, err)
		}
		body = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errorFromResponse(method, path, resp)
	}
	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode %s %s: %w", method, path, err)
	}
	return nil
}

func errorFromResponse(method string, path string, resp *http.Response) error {
	var body models.ErrorResponse
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return Error{
		Method:     method,
		Path:       path,
		StatusCode: resp.StatusCode,
		Code:       body.Code,
		Message:    body.Error,
	}
}

func normalizeBaseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("api base URL is required")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse api base URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("api base URL must use http or https, got %q", u.Scheme)
	}
	if u.Host == "" {
		return "", errors.New("api base URL must include a host")
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return "", errors.New("api base URL must not include query parameters or a fragment")
	}

	return strings.TrimRight(u.String(), "/"), nil
}
