package http_cache

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

func NewCacheTransport(originalTransport http.RoundTripper, seconds int) *CacheTransport {
	transport := &CacheTransport{
		data:              make(map[string]string),
		originalTransport: originalTransport,
	}

	cacheClearJob(seconds, transport)

	return transport
}

func cacheClearJob(seconds int, transport *CacheTransport) {

	interval := time.Duration(seconds) * time.Second
	ticker := time.NewTicker(interval)

	go func() {
		for {

			select {
			case <-ticker.C:
				transport.Clear()
			}
		}
	}()
}

func cacheKey(r *http.Request) string {
	return r.URL.String()
}

type CacheTransport struct {
	data              map[string]string
	mu                sync.RWMutex
	originalTransport http.RoundTripper
}

func (c *CacheTransport) Set(r *http.Request, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[cacheKey(r)] = value
}

func (c *CacheTransport) Get(r *http.Request) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.data[cacheKey(r)]; ok {
		return val, nil
	}

	return "", errors.New("key not found in cache")
}

// Here is the main functionality
func (c *CacheTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	// Check if we have the response cached..
	// If yes, we don't have to hit the server
	// We just return it as is from the cache store.
	if val, err := c.Get(r); err == nil {
		return cachedResponse([]byte(val), r)
	}

	// Ok, we don't have the response cached, the store was probably cleared.
	// Make the request to the server.
	resp, err := c.originalTransport.RoundTrip(r)

	if err != nil {
		return nil, err
	}

	// Get the body of the response so we can save it in the cache for the next request.
	buf, err := httputil.DumpResponse(resp, true)

	if err != nil {
		return nil, err
	}

	// Saving it to the cache store
	c.Set(r, string(buf))

	return resp, nil
}

func (c *CacheTransport) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]string)
	return nil
}

func cachedResponse(b []byte, r *http.Request) (*http.Response, error) {
	buf := bytes.NewBuffer(b)
	return http.ReadResponse(bufio.NewReader(buf), r)
}
