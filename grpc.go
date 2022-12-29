package tools

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func GrpcConnect(addr string, tlsEnabled bool) (*grpc.ClientConn, error) {

	if tlsEnabled {
		creds := credentials.NewTLS(&tls.Config{})
		return grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	}
	return grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

}
