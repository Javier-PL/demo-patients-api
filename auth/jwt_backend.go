package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Ramso-dev/log"
)

var Log log.Logger

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var databasename = "ccl"
var collectionname = "ccl.users"

var authBackendInstance *JWTAuthenticationBackend = nil

//InitJWTAuthenticationBackend retrieves the Keys and returns them in an authBackendInstance struct
func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

//getTokenRemainingValidity checks if the token has expired
func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func getPrivateKey() *rsa.PrivateKey {

	PrivateKeyString := os.Getenv("PRIVATE_KEY")

	if PrivateKeyString == "" {
		Log.Error("PRIVATE_KEY empty")
		return nil
	}

	r := strings.NewReader(PrivateKeyString)
	pemBytes, err := ioutil.ReadAll(r)
	if err != nil {
		Log.Error(err)
		return nil
	}

	data, _ := pem.Decode([]byte(pemBytes))
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		Log.Error(err)
		return nil
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {

	PublicKeyString := os.Getenv("PUBLIC_KEY")

	if PublicKeyString == "" {

		Log.Error("PUBLIC_KEY empty")
		return nil

	}

	r := strings.NewReader(PublicKeyString)
	pemBytes, err := ioutil.ReadAll(r)
	if err != nil {
		Log.Error(err)
		return nil
	}

	data, _ := pem.Decode([]byte(pemBytes))
	publicKeyImported, err := x509.ParsePKCS1PublicKey(data.Bytes)

	if err != nil {
		Log.Error(err)
		return nil
	}

	/*
		rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

		if !ok {
			Log.Error(err)
			return nil
		}*/

	return publicKeyImported
}
