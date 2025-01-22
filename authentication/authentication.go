package authentication

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var publicSigningKey *ecdsa.PublicKey
var privateSigningKey *ecdsa.PrivateKey

func init() {
	var err error

	publicSigningKey, privateSigningKey, err = getPublicPrivateKey("./private-key.pem")

	if err != nil {
		panic(err)
	}
}

func getPublicPrivateKey(filepath string) (*ecdsa.PublicKey, *ecdsa.PrivateKey, error) {
	// first we will try to read the file
	fileData, err := os.ReadFile(filepath)

	if err != nil {
		return nil, nil, err
	}

	// now we will decode it
	pemData, _ := pem.Decode(fileData)

	if pemData == nil {
		return nil, nil, errors.New("failed to decode PEM block")
	}

	// now we will check the headers
	if pemData.Type != "EC PRIVATE KEY" && pemData.Type != "PRIVATE KEY" {
		return nil, nil, errors.New("invalid PEM type")
	}

	// Parse the private key
	var privateKey *ecdsa.PrivateKey

	// First try PKCS#8 format
	key, err := x509.ParsePKCS8PrivateKey(pemData.Bytes)
	if err == nil {
		var ok bool
		privateKey, ok = key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, nil, fmt.Errorf("not an ECDSA private key")
		}
		return &privateKey.PublicKey, privateKey, nil
	}

	// If PKCS#8 fails, try SEC1 format
	privateKey, err = x509.ParseECPrivateKey(pemData.Bytes)
	if err != nil {
		// If both parsing attempts fail, return detailed error
		return nil, nil, fmt.Errorf("failed to parse private key (tried PKCS#8 and SEC1 formats). Make sure the key uses a supported curve (P-224, P-256, P-384, P-521)")
	}

	// Verify the curve is supported
	curve := privateKey.Curve
	switch curve {
	case elliptic.P224(), elliptic.P256(), elliptic.P384(), elliptic.P521():
		// These curves are supported
	default:
		return nil, nil, fmt.Errorf("unsupported elliptic curve: %T", curve)
	}

	return &privateKey.PublicKey, privateKey, nil
}

func CreateToken(username string) (string, error) {
	claims := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		jwt.MapClaims{
			"sub": username,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour).Unix(),
			"iss": "pluto-web",
		},
	)

	jwtString, err := claims.SignedString(privateSigningKey)

	return jwtString, err
}

func VerifyToken(jwtString string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(jwtString, func(t *jwt.Token) (interface{}, error) { return publicSigningKey, nil })

	return jwtToken, err
}
