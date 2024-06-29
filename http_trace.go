package log

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type (
	trace struct {
		statusCode int
		reqBody    any
		respHeader any
		respBody   any
		timeStart  time.Time
	}
)

func NewTrace() trace {
	return trace{
		timeStart: time.Now(),
	}
}

func (t *trace) Save(ctx context.Context, req *http.Request, resp *http.Response) {
	if req == nil {
		return
	}

	// Reading request body
	if err := json.NewDecoder(req.Body).Decode(&t.reqBody); err != nil {
		t.reqBody, _ = io.ReadAll(req.Body)
	}

	// Reading response body
	if resp != nil {
		t.respHeader = resp.Header
		t.statusCode = resp.StatusCode
		if err := json.NewDecoder(resp.Body).Decode(&t.respBody); err != nil {
			if t.respBody, err = io.ReadAll(resp.Body); err != nil {
				Context(ctx).Errorf("failed reading http response body, %v", err)
			}
		}
	}

	globalLogger.LogAttrs(ctx, LevelTrace, "",
		slog.String("caller", GetCaller("", 3)),
		slog.String(processID, Context(ctx).ProcessID()),
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.Int("statusCode", t.statusCode),
		slog.Int64("requestDuration", time.Since(t.timeStart).Milliseconds()),
		slog.Any("requestHeader", req.Header),
		slog.Any("requestBody", t.reqBody),
		slog.Any("responseHeader", t.respHeader),
		slog.Any("responseBody", t.respBody),
	)
}
