package dpop

import (
	"os"
	"testing"
	"context"
	"time"
	"github.com/rs/zerolog"
)

var accessToken = "myaccesstoken"

var logger = zerolog.New(os.Stdout).
				Level(zerolog.DebugLevel).
				With().
				Str("component", "testgocore.dpop").
				Logger()

func TestGoCore_DPoP(t *testing.T){
	var dpop DPoP
	var keyPEM *KeyPEM

	_, cancel := context.WithTimeout(context.Background(),
									time.Duration( 30 ) * time.Second)
	defer cancel()

	// ------ CLIENT SIDE ------	
	// Step 1: client generates the private and public keys.
	dpop = *NewDPoP(&logger)

	keyPEM, err := dpop.CreateKeys()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("keyPEM: %v", keyPEM)
	// ------ CLIENT SIDE ------	

	// ------ CLIENT SIDE ------
	// Step 2: client creates a DPoP JWT token from the authorization server (requesting an access token).
	token, err := dpop.CreateTokenDpopNoAuth("POST",
											 "https://auth.example.com/token", 
											 keyPEM.privPEM,
											 keyPEM.pubPEM,)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token: %v", token)
	// ------ CLIENT SIDE ------

	// ------ SERVER SIDE ------
	bearerToken, err := AuthorizationServer("testuser", 
											token,
											"POST",
											"https://auth.example.com/token")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("BearerToken: %v", bearerToken)
	// ------ SERVER SIDE ------

	// ------ CLIENT SIDE ------
	tokenPop, err := dpop.CreateTokenDPopWithAccessToken("GET",
									"https://api.example.com/orders/123", 
									bearerToken.Token, 
									keyPEM.privPEM,
									keyPEM.pubPEM,)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tokenPop: %v", tokenPop)
	// ------ CLIENT SIDE ------

		// ------ SERVER SIDE ------
	err = AuthorizationServerTokenDPopValidation(tokenPop,
												bearerToken.Token,
												"GET",
												"https://api.example.com/orders/123",)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("BearerToken: %v", bearerToken)
	// ------ SERVER SIDE ------

}