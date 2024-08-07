package server

import (
	"fmt"
	"net"
	"os"

	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RunGRPCServer(registerServer func(server *grpc.Server)) {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	grpcEndpoint := fmt.Sprintf(":%s", httpPort)
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	grpcLogrus.ReplaceGrpcLogger(logrusEntry)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcTags.UnaryServerInterceptor(grpcTags.WithFieldExtractor(grpcTags.CodeGenRequestFieldExtractor)),
			grpcLogrus.UnaryServerInterceptor(logrusEntry),
		),
		grpc.ChainStreamInterceptor(
			grpcTags.StreamServerInterceptor(grpcTags.WithFieldExtractor(grpcTags.CodeGenRequestFieldExtractor)),
			grpcLogrus.StreamServerInterceptor(logrusEntry),
		),
	)
	registerServer(grpcServer)

	listener, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.WithField("grpcEndpoint", grpcEndpoint).Info("Starting gRPC listener")
	logrus.Fatal(grpcServer.Serve(listener))
}
