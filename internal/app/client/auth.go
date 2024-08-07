package client

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"google.golang.org/grpc/credentials"
)

type metaServerToken struct {
	serviceURL string
}

func newMetaServerToken(grpcAddress string) credentials.PerRPCCredentials {
	return metaServerToken{serviceURL: "https://" + strings.Split(grpcAddress, ":")[0]}
}

func (m metaServerToken) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", m.serviceURL)

	idToken, err := metadata.GetWithContext(ctx, tokenURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to request id token for gRPC: %w", err)
	}

	return map[string]string{"authorization": "Bearer " + idToken}, nil
}

func (m metaServerToken) RequireTransportSecurity() bool {
	return true
}
