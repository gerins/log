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
		Time           time.Time   `json:"time"`
		Method         string      `json:"method"`
		Url            string      `json:"-"`
		StatusCode     int         `json:"statusCode"`
		Duration       int64       `json:"duration"`
		ReqHeader      any         `json:"requestHeader"`
		ReqBody        any         `json:"requestBody"`
		RespHeader     http.Header `json:"responseHeader"`
		RespBody       any         `json:"responseBody"`
		RawRespBody    []byte      `json:"-"`
		addToExtraData bool        `json:"-"`
	}
)

func NewTrace(method, url string, reqHeader, reqBody any, addToExtraData bool) trace {
	return trace{
		Time:           time.Now(),
		Method:         method,
		Url:            url,
		ReqHeader:      reqHeader,
		ReqBody:        reqBody,
		addToExtraData: addToExtraData,
	}
}

func (t *trace) Save(ctx context.Context, resp *http.Response) {
	if err := json.Unmarshal(t.RawRespBody, &t.RespBody); err != nil {
		t.RespBody = string(t.RawRespBody)
	}

	// Reading response body
	if resp != nil {
		t.RespHeader = resp.Header
		t.StatusCode = resp.StatusCode
	}

	t.Duration = time.Since(t.Time).Milliseconds()

	if t.addToExtraData {
		Context(ctx).ExtraData[t.Url] = t
		return
	}

	globalLogger.LogAttrs(ctx, LevelTrace, "",
		slog.String("caller", GetCaller("", 4)),
		slog.String(processID, Context(ctx).ProcessID()),
		slog.String("method", t.Method),
		slog.String("url", t.Url),
		slog.Int("statusCode", t.StatusCode),
		slog.Int64("requestDuration", t.Duration),
		slog.Any("requestHeader", t.ReqHeader),
		slog.Any("requestBody", t.ReqBody),
		slog.Any("responseHeader", t.RespHeader),
		slog.Any("responseBody", t.RespBody),
	)
}
