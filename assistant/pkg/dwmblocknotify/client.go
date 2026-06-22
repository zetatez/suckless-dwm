package dwmblocknotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	notifyURL = "http://127.0.0.1:8765/notify"
	timeout   = 3 * time.Second
)

type Client struct {
	HTTP *http.Client
}

var (
	client     *Client
	clientOnce sync.Once
)

// getClient returns the singleton Client.
func getClient() *Client {
	clientOnce.Do(func() {
		client = &Client{HTTP: &http.Client{Timeout: timeout}}
	})
	return client
}

type notifyBody struct {
	Msg        string `json:"msg"`
	TTLSeconds int    `json:"ttl_seconds"`
}

// POST asynchronously enqueues a message (shown after previous messages expire).
func POST(msg string, ttl time.Duration) {
	go sendJSON(http.MethodPost, msg, ttl)
}

// PUT asynchronously preempts the current message; the preempted one is re-queued at the head.
func PUT(msg string, ttl time.Duration) {
	go sendJSON(http.MethodPut, msg, ttl)
}

// DELETE asynchronously clears the current message and the entire queue.
func DELETE() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, notifyURL, nil)
		if err != nil {
			logErr(http.MethodDelete, err)
			return
		}
		do(req)
	}()
}

func sendJSON(method, msg string, ttl time.Duration) {
	buf, err := json.Marshal(notifyBody{Msg: msg, TTLSeconds: int(ttl.Seconds())})
	if err != nil {
		logErr(method, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, notifyURL, bytes.NewReader(buf))
	if err != nil {
		logErr(method, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	do(req)
}

func do(req *http.Request) {
	resp, err := getClient().HTTP.Do(req)
	if err != nil {
		logErr(req.Method, err)
		return
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
}

func logErr(method string, err error) {
	fmt.Fprintf(os.Stderr, "dwmblocknotify %s failed: %v\n", method, err)
}
