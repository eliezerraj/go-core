package httpclient

import (
//	"bytes"
//	"context"
//	"fmt"
	"net/http"
//	"net/http/httptrace"
	"sync"
	"time"

	"github.com/eliezerraj/go-core/v3/logger"
)

type HttpConfig struct {
	Timeout             time.Duration
	KeepAlive           time.Duration
	IdleConnTimeout     time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	ServiceName         string
}

type IHTTPClient interface {
	Do(*http.Request) (*http.Response, error)
	CloseIdleConnections()
}

type Client struct {
	*http.Client
	mu sync.Mutex
}

type Requester struct {
	Scheme  string
	Host    string
	Path    string
	Headers map[string]string
	Method  string
	Body    []byte
	Logger  logger.ILogger
}

// New creates a new HTTP client with the provided configuration.
func NewHttpClient(cfg *HttpConfig) *Client {
	transport := &http.Transport{
		//DialTLSContext:      customDialTLSContext(cfg),
		MaxIdleConns:        cfg.MaxIdleConns,
		MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
		MaxConnsPerHost:     cfg.MaxConnsPerHost,
		IdleConnTimeout:     cfg.IdleConnTimeout,
		DisableKeepAlives:   false,
		ForceAttemptHTTP2:   false,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	return &Client{
		httpClient, sync.Mutex{},
	}
}

func NewRequester(host, path, method string, headers map[string]string, body []byte, logger logger.ILogger) Requester {
	return Requester{
		Host:    host,
		Path:    path,
		Headers: headers,
		Method:  method,
		Body:    body,
		Logger:  logger,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Client.Do(req)
}