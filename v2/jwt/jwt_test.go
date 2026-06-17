package jwt

import (
	"context"
	"testing"
)

var token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik5QMkxpdzVUbEZtZFNWbUUtdVYzYiJ9.eyJpc3MiOiJodHRwczovL2Rldi1temRxeTRzOC51cy5hdXRoMC5jb20vIiwic3ViIjoia2xNbXp1bGVycmtueVlRVFBlZlgyOFpscVFOVjBaTUJAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vZGV2LW16ZHF5NHM4LnVzLmF1dGgwLmNvbS9hcGkvdjIvIiwiaWF0IjoxNzgxMzgyOTI3LCJleHAiOjE3ODE0NjkzMjcsInNjb3BlIjoicmVhZDp1c2VycyB1cGRhdGU6dXNlcnMgZGVsZXRlOnVzZXJzIGNyZWF0ZTp1c2VycyIsImd0eSI6ImNsaWVudC1jcmVkZW50aWFscyIsImF6cCI6ImtsTW16dWxlcnJrbnlZUVRQZWZYMjhabHFRTlYwWk1CIn0.H0flNpM4CWl-RtedMjXxmgjtoGsbBSmG_3-0u8YnPO858rhOMwiI4MAUEDRdX-VJV0eddcFKi9F2PvnKyud5SbWJZQjXQNsQNVZ2aFoEg4gPJEA3rREfImgHqhTNoh6X3RzFtz9cZX04nEAwtZBjNmcSqrBlb4zCxvnJiIlewuHQrv71EJVxrMbTQLDdOV14vk2VJtLH9O_R17xZrlWjzXyyp8ISQok6HttuxcajKGpP6D2Gx87p4DbZ70St7z2EtZfJVFssLvF8GLMShWDkTZm9i4T-F35-KVOHCD-RnJJWlyJFvQZA6WFZkPvaXwyliM7IH3-Fcg9IF1mJg2XzKw"
var publicKeyN = "wQQaeRjWm9563jb61cminNZefrMlzxaEHg0ZAHKukCJj41Ewb4JuhMmwSUiiCLQ3qP3BURBIf400n1ZBD3Lylpn2nKNaBOdD0bd0uWXtNHfYQs9bacCjM24f61heGMYRL1YHgIrHvTTUiJBG5Y-TJzDIjmPhfJtTLU6jpvDzj4yqej5SRx3MOT9WuzIJDkkvxid1bLjOZ9ZkjnYIHojkI7jCItLFtBS81OiJdCCC1OC_uUoyCQDKrxUHVb6eZK9XhNcgbJA7hRE1BtrEgiVKKbfcZQbcj4eciKyMpYmH9ZBjeP-RSDiHFdyTdAbsl_9f4Yn-WuiC_gAjQVgxzEAnxQ"
var publicKeyE = "AQAB"

var token2 = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik5QMkxpdzVUbEZtZFNWbUUtdVYzYiJ9.eyJpc3MiOiJodHRwczovL2Rldi1temRxeTRzOC51cy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NmEyZGEwYTg3YjQ0NjA2YWJmNGZkMmY3IiwiYXVkIjpbImh0dHBzOi8vZGV2LW16ZHF5NHM4LnVzLmF1dGgwLmNvbS9hcGkvdjIvIiwiaHR0cHM6Ly9kZXYtbXpkcXk0czgudXMuYXV0aDAuY29tL3VzZXJpbmZvIl0sImlhdCI6MTc4MTQ4MjI2NywiZXhwIjoxNzgxNTY4NjY3LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIHJlYWQ6Y3VycmVudF91c2VyIHVwZGF0ZTpjdXJyZW50X3VzZXJfbWV0YWRhdGEgZGVsZXRlOmN1cnJlbnRfdXNlcl9tZXRhZGF0YSBjcmVhdGU6Y3VycmVudF91c2VyX21ldGFkYXRhIGNyZWF0ZTpjdXJyZW50X3VzZXJfZGV2aWNlX2NyZWRlbnRpYWxzIGRlbGV0ZTpjdXJyZW50X3VzZXJfZGV2aWNlX2NyZWRlbnRpYWxzIHVwZGF0ZTpjdXJyZW50X3VzZXJfaWRlbnRpdGllcyIsImd0eSI6InBhc3N3b3JkIiwiYXpwIjoiUzgyMDRSUzRKWnlHcndNUm9jUEp1eml5c3R6eWhOM3EifQ.bWjurslZIaO7G2FU13e7gDgMLWU4fXwOihQnWxqijDDUtQPxTMkZtxIatYvWctSAp3N4vL4qx7-izyoABHeOj1a0sW_OxToOcIEBr67xbIptendQjJdW0XwsnA_wtElQavHR_l7Xtf-XuM4E9E_3lyB8tTSGejpwFeLD0G4ITgIawJ98IWR431CfkPyQzLANlaWafljPtPQduSKi9NNoaaVlAJi8Ewq_E343damUvbzEcgzB2HP4FnUcJq3eoCvDSNufgsFky6OK5c4TgwUQF71gdpJP7pXJyudozqopaoxE9hnIXsS7Qr_12rm2WJHKOQpr-Na_LvHDaocp-0wLAA"
var publicKeyN2 = "wQQaeRjWm9563jb61cminNZefrMlzxaEHg0ZAHKukCJj41Ewb4JuhMmwSUiiCLQ3qP3BURBIf400n1ZBD3Lylpn2nKNaBOdD0bd0uWXtNHfYQs9bacCjM24f61heGMYRL1YHgIrHvTTUiJBG5Y-TJzDIjmPhfJtTLU6jpvDzj4yqej5SRx3MOT9WuzIJDkkvxid1bLjOZ9ZkjnYIHojkI7jCItLFtBS81OiJdCCC1OC_uUoyCQDKrxUHVb6eZK9XhNcgbJA7hRE1BtrEgiVKKbfcZQbcj4eciKyMpYmH9ZBjeP-RSDiHFdyTdAbsl_9f4Yn-WuiC_gAjQVgxzEAnxQ"
var publicKeyE2 = "AQAB"

func TestJWTValidator(t *testing.T) {
	jwtValidator := (&JWTValidator{}).NewJWTValidatorFromJWK(publicKeyN, publicKeyE)
	if jwtValidator == nil {
		t.Fatalf("FAILED to create JWT validator from JWK")
	}

	valid, err := jwtValidator.ValidateToken(context.Background(), token)
	if err != nil {
		t.Errorf("FAILED to validate JWT token: %s", err)
	}
	if !valid {
		t.Errorf("JWT token is not valid")
	} else {
		t.Logf("JWT token is valid")
	}
}

func TestJWTValidator2(t *testing.T) {
	jwtValidator := (&JWTValidator{}).NewJWTValidatorFromJWK(publicKeyN2, publicKeyE2)
	if jwtValidator == nil {
		t.Fatalf("FAILED to create JWT validator from JWK")
	}

	valid, err := jwtValidator.ValidateToken(context.Background(), token2)
	if err != nil {
		t.Errorf("FAILED to validate JWT token: %s", err)
	}
	if !valid {
		t.Errorf("JWT token is not valid")
	} else {
		t.Logf("JWT token is valid")
	}
}

func TestJWTValidatorExpiredToken(t *testing.T) {
	jwtValidator := (&JWTValidator{}).NewJWTValidatorFromJWK(publicKeyN2, publicKeyE2)
	if jwtValidator == nil {
		t.Fatalf("FAILED to create JWT validator from JWK")
	}

	valid, err := jwtValidator.ValidateToken(context.Background(), token2)
	if err != nil {
		t.Fatalf("token is validation error: %s", err)
	}
	if valid {
		t.Logf("token IS NOT EXPIRED")
	} else {		
		t.Logf("token IS EXPIRED")
	}
}
