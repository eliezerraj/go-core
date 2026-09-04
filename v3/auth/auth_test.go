package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetJwksUrl(t *testing.T) {
	ctx := context.Background()

	jwksUrl := "http://localhost:7100/v1/.well-known/jwks.json"

	authService := NewAuthService(jwksUrl, false, "Authorization", 5*time.Second)

	err := authService.GetJwksUrl(ctx)
	if err != nil {
		t.Fatalf("failed to get JWKS URL: %v", err)
	}
	assert.NoError(t, err)

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImRlZmF1bHQta2lkIiwidHlwIjoiSldUIn0.eyJjbGllbnRfaWQiOiJ0ZXN0Z2UiLCJzY29wZSI6InRlc3RlOnJlYWQiLCJpc3MiOiJnby1hdXRob3JpemVyLXYyIiwic3ViIjoidGVzdGdlIiwiYXVkIjpbImF1ZC10ZXN0ZSJdLCJleHAiOjE3ODg0ODA4MzYsIm5iZiI6MTc4ODQ3NzIzNiwiaWF0IjoxNzg4NDc3MjM2LCJqdGkiOiJlODA1NmUxYy0zNjg4LTRhNDEtODcyNi1hYTE3MmI0NWFjYmUifQ.PryjtTJqDvGMJO8NWKrRm85DkospSFOjmuo5V2cZWAKl7qjnQ_AJWOB7NApVX0iU-m-3UwFX-fIzendC_DtPKEp3vcBJ67otpuXDVjf_hjoJIU5P78b-2qObxqJpNsaxnECZrVIj12yG3hW37_oN843DAaM-Wf5FkGyo6wf1HHoOdSVxg6nfqSA9tbB9NvEfyMb3mYotuwMBaHyWfTaoFLuwDuCeSCeQEJCIYYVi5dw7qy3CT89cS78RxiDHHV57jwjjEc3wzkxlEPMGVfcrnF357V4OVjmmQ4x-TO-ZhKerhkNnBZeFQLT2SSVHbEUPX282vs4_FcjiwEbENhD36lnbUCbDRuKb7T3bYBlrbHHxHuC3wj72AsnH7rsPdzzvAqmTh3Gyw-8MNG1amTJB4eJ8ksNTu6IFmX8WWwjMWtM5mok2n8xX_ez-eI5Q8bn3Hvq-NmWkc7GnBB1sOFKj-pAA8iCoqqEcqNMaQSNPWSsK9Zfwy5BrETcXrWjJoB7j"
	claims, err := authService.VerifyToken(ctx, token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)

}