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
				Level(zerolog.WarnLevel).
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
	t.Logf("==================================================")
	t.Logf("=========== 1 - CLIENT SIDE: Generating keys =========\n")

	dpop = *NewDPoP(&logger)

	keyPEM, err := dpop.CreateKeys()
	if err != nil {
		t.Errorf("err : %s", err)
	}
	t.Logf("keyPEM: %v", keyPEM)
	// ------ CLIENT SIDE ------	

	// ------ CLIENT SIDE ------
	t.Logf("=========== 2 - CLIENT SIDE: Creating Client DPoP (SIGNED RSA) =========\n")

	// Step 2: client creates a DPoP JWT token from the authorization server (requesting an access token).
	client_token, err := dpop.CreateTokenDpopNoAuth("POST",
											 "https://auth.example.com/token", 
											 keyPEM.privPEM,
											 keyPEM.pubPEM,)
	if err != nil {
		t.Errorf("err : %s", err)
	}
	t.Logf("-----> client_token : %v", client_token)
	// ------ CLIENT SIDE ------
	
	t.Logf("")
	
	// ------ SERVER SIDE ------
	t.Logf("....................................................")		
	t.Logf("... SERVER SIDE: Requesting Bearer Token (dpop embedded) ...\n")

	serverBearerToken, err := AuthorizationServer("testuser", 
											client_token,
											"POST",
											"https://auth.example.com/token")
	if err != nil {
		t.Errorf("err : %s", err)
	}
	t.Logf("----> serverBearerToken: %v", serverBearerToken)
	// ------ SERVER SIDE ------

	t.Logf("")

	// ------ CLIENT SIDE ------
	t.Logf("======================================================================")
	t.Logf("======== 3 - CLIENT SIDE: Creating DPoP token with access token =====\n")
	clientTokenDPop, err := dpop.CreateTokenDPopWithAccessToken("GET",
									"https://api.example.com/orders/123", 
									serverBearerToken.Token, 
									keyPEM.privPEM,
									keyPEM.pubPEM,)
	if err != nil {
		t.Errorf("err : %s", err)
	}
	t.Logf("----> clientTokenDPop: %v", clientTokenDPop)
	// ------ CLIENT SIDE ------

	t.Logf("")

	// ------ SERVER SIDE ------
	t.Logf("....................................................")		
	t.Logf(".... 4 - SERVER SIDE: Validating DPoP token with access token .....")
	err = AuthorizationServerTokenDPopValidation(clientTokenDPop,
												serverBearerToken.Token,
												"GET",
												"https://api.example.com/orders/123",)
	if err != nil {
		t.Errorf("err : %s", err)
	} else {
		t.Logf("----> DPoP token with access token validated successfully")
	}
	
	t.Logf("....................................................")		
	t.Logf(".... 5 - SERVER SIDE: Validating DPoP token with access token .....")
	err = AuthorizationServerTokenDPopValidation(clientTokenDPop,
												serverBearerToken.Token,
												"GET",
												"https://api.example.com/orders/456",)
	if err != nil {
		t.Errorf("err : %s", err)
	} else {
		t.Logf("----> DPoP token with access token validated successfully")
	}

		// ------ CLIENT SIDE ------
	t.Logf("======================================================================")
	t.Logf("======== 6 - CLIENT SIDE: Creating DPoP token with access token =====\n")
	clientTokenDPop, err = dpop.CreateTokenDPopWithAccessToken("GET",
									"https://api.example.com/orders/456", 
									serverBearerToken.Token, 
									keyPEM.privPEM,
									keyPEM.pubPEM,)
	if err != nil {
		t.Errorf("err : %s", err)
	} else {
		t.Logf("----> DPoP token with access token validated successfully")
	}

	t.Logf("----> clientTokenDPop: %v", clientTokenDPop)
	// ------ CLIENT SIDE ------

	t.Logf("")

	// ------ SERVER SIDE ------
	t.Logf("....................................................")		
	t.Logf(".... 7 - SERVER SIDE: Validating DPoP token with access token .....")
	err = AuthorizationServerTokenDPopValidation(clientTokenDPop,
												serverBearerToken.Token,
												"GET",
												"https://api.example.com/orders/456",)
	if err != nil {
		t.Errorf("err : %s", err)
	} else {
		t.Logf("----> DPoP token with access token validated successfully")
	}
	
	// ------ SERVER SIDE ------
}