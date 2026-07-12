
Step 1: Client generates a key pair

## Flow 1 - Requesting an Access Token

Step 2: Client request an access token DPoP proof

DPoP structure

```
JWT Header
-----------
alg = ES256
typ = dpop+jwt
jwk = Public Key

JWT Payload
------------
htm = POST
htu = https://auth.example.com/token
iat = ...
jti = ...

Signature
----------
Signed with Private Key
```

Diagram

```
Client (PK + PB)
    │
    │ POST /token
    │
    │ DPoP: <JWT>
        {
        "htm":"POST",
        "htu":"https://auth.example.com/token",
        "iat":12345,
        "jti":"abc"
        }
    │
    ▼
Authorization Server
   │
   │ verifies proof
   │
   │ creates NEW JWT + JKT
   │ (signed by AUTH SERVER)
   ▼
Access Token
```

Step 3: Authorization Server

The Authorization Server can only validate:

    ✅ JWT signature
    ✅ htm
    ✅ htu
    ✅ iat
    ✅ jti

There is no ath claim because the access token doesn't exist yet.

The Authorization Server then issues: 

    Access Token

```
{
   "sub": "alice",
   "exp": 1751888000,
   "cnf": {
       "jkt": "SHA256(PublicKeyA)"
   }
}
```

The jkt value is the thumbprint (fingerprint) of the public key.

## Flow 2 - Calling an API

Now the client already has:

    Access Token + Private Key

The Authorization Server issued this access token:

```json
{
    "sub":"alice",
    "scope":"orders.read",
    "exp":1751888000,

    "cnf":{
        "jkt":"ABC123XYZ"
    }
}
```

Diagram

```
Client (PK + PB + Access Token)
    │
    │ GET /order
    │
    │ DPoP: NEW <JWT>
        {
            "alg":"ES256",
            "typ":"dpop+jwt",
            "jwk":{
                "kty":"EC",
                "crv":"P-256",
                "x":".....",
                "y":"....."
            }
        }
    │
    ▼
Authorization Server
   │
   │ verifies proof
   │
   │ creates NEW JWT + JKT
   │ (signed by AUTH SERVER)
   ▼
Access Token
```


The request becomes:

    GET /orders

    Authorization: DPoP eyJhbGc...accessToken...

    DPoP: eyJhbGc...proof...

Note: A new proof is created for every request.



Second DPoP JWT (with access code)
```json
{
   "htm":"GET",
   "htu":"https://api.example.com/orders",
   "iat":12345,
   "jti":"xyz",
   "ath":"..."
}
```

Step 3: Client sends the DPoP proof

    POST /token

    DPoP: eyJhbGciOiJFUzI1NiIs...

Step 4: Server validates the proof

```
                 CLIENT

Generate Key Pair
-----------------
Private Key
Public Key

        │
        │ Build JWT
        │
        ▼

JWT Header
-----------
Public Key

JWT Payload
------------
POST /token
iat
jti

        │
        │ Sign with Private Key
        ▼

DPoP JWT
        │
        ▼

=============== HTTP ===============

POST /token
DPoP: eyJhbGc...

====================================

        │
        ▼

          AUTHORIZATION SERVER

Extract Public Key
        │
        ▼
Verify Signature
        │
        ▼
Issue Access Token
(bound to that public key)
```