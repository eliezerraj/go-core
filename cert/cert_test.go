package cert

import (
	"io/ioutil"
	"testing"
)

var public_key = ""

func TestCore_ParsePemToRSAPriv(t *testing.T){
	var certCore CertCore

	private_key, err := ioutil.ReadFile("client_priv.key")
	if err != nil {
		t.Errorf("Failed to read private key: %v", err)
	}
	str := string(private_key)
	_, err = certCore.ParsePemToRSAPriv(&str)
	if err != nil {
		t.Errorf("failed to ParsePemToRSAPriv : %s", err)
	}

	public_key, err := ioutil.ReadFile("client_pub.key")
	if err != nil {
		t.Errorf("Failed to read private key: %v", err)
	}
	str = string(public_key)
	_, err = certCore.ParsePemToRSAPub(&str)
	if err != nil {
		t.Errorf("failed to ParsePemToRSAPub : %s", err)
	}
}