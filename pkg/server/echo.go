package server

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/takama/grpc/contracts/echo"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type echoServer struct {
	log *zap.Logger
}

// Echo the content of the request
func (es echoServer) Ping(ctx context.Context, in *echo.Request) (*echo.Response, error) {
	es.log.Info("ping:", zap.String("Handling echo request", fmt.Sprintf("[%v] with context %v", in, ctx)))
	hostname, err := os.Hostname()
	if err != nil {
		es.log.Error("ping: Unable to get hostname", zap.Error(err))
		hostname = ""
	}
	err = grpc.SendHeader(ctx, metadata.Pairs("hostname", hostname))
	if err != nil {
		es.log.Error("ping: send header", zap.Error(err))
	}
	time.Sleep(time.Millisecond * 200)
	return &echo.Response{Content: in.Content}, nil
}

// Reverse the content of the request
func (es echoServer) Reverse(ctx context.Context, in *echo.Request) (*echo.Response, error) {
	es.log.Info("reverse:", zap.String("Handling reverse request", fmt.Sprintf("[%v] with context %v", in, ctx)))
	hostname, err := os.Hostname()
	if err != nil {
		es.log.Error("reverse: Unable to get hostname", zap.Error(err))
		hostname = ""
	}
	err = grpc.SendHeader(ctx, metadata.Pairs("hostname", hostname))
	if err != nil {
		es.log.Error("ping: send header", zap.Error(err))
	}
	return &echo.Response{Content: reverse(in.Content)}, nil
}

func reverse(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < len(runes)/2; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
