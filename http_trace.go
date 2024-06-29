package log

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type (
	trace struct {
		method, url        string
		reqHeader, reqBody any
		respHeader         http.Header
		RespBody           []byte
		statusCode         int
		timeStart          time.Time
	}
)

func NewTrace(method, url string, reqHeader, reqBody any) trace {
	return trace{
		method:    method,
		url:       url,
		reqHeader: reqHeader,
		reqBody:   reqBody,
		timeStart: time.Now(),
	}
}

func (t *trace) Save(ctx context.Context, resp *http.Response) {
	var respModel any
	if err := json.Unmarshal(t.RespBody, &respModel); err != nil {
		respModel = string(t.RespBody)
	}

	// Reading response body
	if resp != nil {
		t.respHeader = resp.Header
		t.statusCode = resp.StatusCode
	}

	globalLogger.LogAttrs(ctx, LevelTrace, "",
		slog.String("caller", GetCaller("", 3)),
		slog.String(processID, Context(ctx).ProcessID()),
		slog.String("method", t.method),
		slog.String("url", t.url),
		slog.Int("statusCode", t.statusCode),
		slog.Int64("requestDuration", time.Since(t.timeStart).Milliseconds()),
		slog.Any("requestHeader", t.reqHeader),
		slog.Any("requestBody", t.reqBody),
		slog.Any("responseHeader", t.respHeader),
		slog.Any("responseBody", respModel),
	)
}
