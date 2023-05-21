package grpc

import (
	"context"
	"fmt"
	"runtime"

	"github.com/spf13/cast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/gerins/log"
)

var (
	stackSize         = 4 << 10 // 4 KB
	statusCodeSuccess = 200
	statusCodeFailed  = 500
)

func SaveLogRequest() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		requestLog := log.NewRequest()
		ctx = requestLog.SaveToContext(ctx)

		// Get request metadata from context
		if requestMetadata, ok := metadata.FromIncomingContext(ctx); ok {
			requestLog.ReqHeader = requestMetadata

			if processID := requestMetadata.Get("process_id"); len(processID) != 0 {
				requestLog.SetProcessID(processID[0])
			}
			if userID := requestMetadata.Get("user_id"); len(userID) != 0 {
				requestLog.UserID = cast.ToInt(userID[0])
			}
		}

		// Get client IP Address from context
		if client, ok := peer.FromContext(ctx); ok {
			requestLog.IP = client.Addr.String()
		}

		// Recover if panic occur
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				stack := make([]byte, stackSize)
				length := runtime.Stack(stack, false)
				requestLog.Debug(fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length]))
			}

			requestLog.Save()
		}()

		// Proceed request
		response, err := handler(ctx, req)

		requestLog.ReqBody = req
		requestLog.RespBody = response
		requestLog.Method = "GRPC"
		requestLog.URL = info.FullMethod
		requestLog.StatusCode = statusCodeSuccess

		if err != nil {
			requestLog.RespBody = err
			requestLog.StatusCode = statusCodeFailed
		}

		return response, err
	}
}
