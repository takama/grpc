package server

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/takama/grpc/contracts/echo"
	"github.com/takama/grpc/pkg/client"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type echoServer struct {
	cl  *client.Client
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

	md := new(metadata.MD)

	cl := echo.NewEchoClient(es.cl.Connection())
	response, err := cl.Ping(es.cl.Content(), &echo.Request{
		Content: in.Content}, grpc.Header(md))
	if err != nil {
		es.log.Error("reverse ping: do request to other service", zap.Error(err))
		return response, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		es.log.Error("reverse ping: Unable to get hostname", zap.Error(err))
		hostname = ""
	}
	err = grpc.SendHeader(ctx, metadata.Pairs("hostname", hostname))
	if err != nil {
		es.log.Error("reverse ping: send header", zap.Error(err))
	}
	return &echo.Response{Content: fmt.Sprintf("%s from %v", response.Content, md.Get("hostname"))}, nil
}
