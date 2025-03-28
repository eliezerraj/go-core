package gprc_client

import (
	"context"
	"testing"
)

func TestMyServiceClient_GetData(t *testing.T) {
	var testGrpcClient GrpcClient

	hostGrpc := "localhost:50053"

	grpcClient, err  := testGrpcClient.StartGrpcClient(hostGrpc)
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}

	err = grpcClient.TestConnection(context.Background())
	if err != nil {
		t.Fatalf("Failed to connected server: %v", err)
	}

	t.Logf("grpcClient: %v", grpcClient)
}