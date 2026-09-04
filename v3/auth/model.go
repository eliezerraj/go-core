package auth
import(
    "github.com/golang-jwt/jwt/v5"
)

type WellKnownJwks struct{
	Keys	[]JWK   `json:"keys"`
}

type JWK struct {
    KeyID     string `json:"kid"`
    KeyType   string `json:"kty"` // RSA
    Algorithm string `json:"alg"` // RS256
    Use       string `json:"use"` // sig
    N         string `json:"n"`   // Base64URL-encoded Modulus
    E         string `json:"e"`   // Base64URL-encoded Exponent
}

type Claims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}
