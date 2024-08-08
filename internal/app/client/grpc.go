package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/lGreenLightl/fitness-club-service/internal/app/genproto/customer"
	"github.com/lGreenLightl/fitness-club-service/internal/app/genproto/trainer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCustomerClient() (client customer.CustomerServiceClient, close func() error, err error) {
	grpcAddress := os.Getenv("CUSTOMER_GRPC_ADDR")
	if grpcAddress == "" {
		return nil, func() error { return nil }, errors.New("empty environment CUSTOMER_GRPC_ADDR")
	}

	options, err := grpcDialOpts(grpcAddress)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	conn, err := grpc.NewClient(grpcAddress, options...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return customer.NewCustomerServiceClient(conn), conn.Close, nil
}

func WaitForCustomerService(timeout time.Duration) bool {
	return isPortAvailable(timeout, os.Getenv("CUSTOMER_GRPC_ADDR"))
}

func NewTrainerClient() (client trainer.TrainerServiceClient, close func() error, err error) {
	grpcAddress := os.Getenv("TRAINER_GRPC_ADDR")
	if grpcAddress == "" {
		return nil, func() error { return nil }, errors.New("empty environment TRAINER_GRPC_ADDR")
	}

	options, err := grpcDialOpts(grpcAddress)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	conn, err := grpc.NewClient(grpcAddress, options...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return trainer.NewTrainerServiceClient(conn), conn.Close, nil
}

func WaitForTrainerService(timeout time.Duration) bool {
	return isPortAvailable(timeout, os.Getenv("TRAINER_GRPC_ADDR"))
}

func isPortAvailable(timeout time.Duration, address string) bool {
	availableChan := make(chan struct{})
	timeoutChan := time.After(timeout)

	go func() {
		for {
			select {
			case <-timeoutChan:
				return
			default:
			}

			_, err := net.Dial("tcp", address)
			if err == nil {
				close(availableChan)
				return
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {
	case <-availableChan:
		return true
	case <-timeoutChan:
		return false
	}
}

func grpcDialOpts(grpcAddress string) ([]grpc.DialOption, error) {
	if ok, _ := strconv.ParseBool(os.Getenv("GRPC_WITH_TLS")); !ok {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("unable load root CA cert: %w", err)
	}
	credentials := credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})

	return []grpc.DialOption{
		grpc.WithTransportCredentials(credentials),
		grpc.WithPerRPCCredentials(newMetaServerToken(grpcAddress)),
	}, nil
}
