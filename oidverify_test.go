package main

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"os"
	"testing"
)

func createJwt(header, payload, signature string) string {
	eHeader := base64.RawURLEncoding.EncodeToString([]byte(header))
	ePayload := base64.RawURLEncoding.EncodeToString([]byte(payload))
	eSignature := base64.RawURLEncoding.EncodeToString([]byte(signature))

	return eHeader + "." + ePayload + "." + eSignature
}

func TestSuccessCreatingJwt(t *testing.T) {
	dHeader := "test_header"
	dPayload := "test_playload"
	dSignature := "test_signature"

	jwt := createJwt(dHeader, dPayload, dSignature)

	decodedJwt, err := NewDecodedJwt(jwt)
	if err != nil {
		t.Error("create fail")
	}

	if decodedJwt.header != dHeader {
		t.Errorf("decodedJwt's header is %s. wanted %s", decodedJwt.header, dHeader)
	}

	if decodedJwt.payload != dPayload {
		t.Errorf("decodedJwt's payload is %s. wanted %s", decodedJwt.payload, dPayload)
	}

	if decodedJwt.signature != dSignature {
		t.Errorf("decodedJwt's signature is %s. wanted %s", decodedJwt.signature, dSignature)
	}
}

// 実際のデータを利用したため、環境変数より取得
// 処理イメージ（goolge の openId connect のみ想定）
func TestVerifySignature(t *testing.T) {
	jwt, err := NewDecodedJwt(os.Getenv("test_jwt"))
	if err != nil {
		t.Error("create jwt fail")
	}
	// 以下の値をハッシュにする
	m := jwt.CalcHash() //hashes

	// public key で利用する e, n は
	// 実際には以下のような場所から取得する
	// https://accounts.google.com/.well-known/openid-configuration
	// https://www.googleapis.com/oauth2/v3/certs
	publicKey, err := NewPublicKey(
		os.Getenv("test_exponent"),
		os.Getenv("test_modulus"),
	)
	if err != nil {
		t.Error("create public_key fail")
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, m, []byte(jwt.signature))
	if err != nil {
		t.Error("don't verify signature")
	}
}
