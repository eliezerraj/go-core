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

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImdvLWF1dGhvcml6ZXItdjIta2V5LTIwMjYtMDkiLCJ0eXAiOiJKV1QifQ.eyJjbGllbnRfaWQiOiJjbGllbnQtdGVzdC0wMSIsInNjb3BlIjoidGVzdGU6cmVhZCIsImlzcyI6ImdvLWF1dGhvcml6ZXItdjIiLCJzdWIiOiJjbGllbnQtdGVzdC0wMSIsImF1ZCI6WyJhdWQtdGVzdGUiXSwiZXhwIjoxNzg4Njk5MDkyLCJuYmYiOjE3ODg2NjMwOTIsImlhdCI6MTc4ODY2MzA5MiwianRpIjoiZWRiNWJjZTktZGRlZi00NTVlLTg0YWYtN2IzMWNkMDI2YjFkIn0.f62lZ8VdCok8GhxAvRO0NXYnapT36CNzKa-k-ST8vDqUvYCAqGC-mYNwROYbcUaimAiicNG8iYNhr2BOBjWflddRJu5-20b_KhbqbqD6Mxa4zjN3lvKE8WxqaUJOQksCdSA8iG-mTOkR56S8--pAZM-XUV8ypw4nKYaK_ffwtst7JpDRqixIuFCSI-UGfkUH72aDA38ST3joMppbRhrblk5F41e3rBL9JREgEhJ32mqgO6nHycbGLN4-S5eJW6RlEFEWUzjMnPVnc3ioDT98Eih5ASBY0NgRl6XET6aqN3cRLOQWPFqOu12G7N-l3y9eQoh4mwSlqO2tGquoR6jcF18oUWbBtJerbeHsiae2bsCsMPKgJZ6ovarAOuvGTICP5sDL8YL6i0PBOlWHwxlzozt-EYZoyX-I7MgnYpCWxE2ucReLXCQPIw0BHTn44Q7lK7NEnVPQduHdK0WA_gwqRzHjlL6pSbsoT9cCzTovmWP03TFUedQDeQPSX-xc2HmW"
	claims, err := authService.VerifyToken(ctx, token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)

	authClientService := NewAuthClientService(func(a *AuthClientService) {
		a.authURL = "http://localhost:7100/v1/login"
		a.refreshURL = "http://localhost:7100/v1/refresh"
		a.clientID = "client-test-01"
		a.clientSecret = "client-secret-test-01"
		a.dryRun = false
		a.refreshInterval = 60
	})
	
	accessToken, err := authClientService.StartAuthenticate(ctx)

	time.Sleep(300 * time.Second)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)

}