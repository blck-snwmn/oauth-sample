package main

import (
	"crypto"
	"crypto/rsa"
	_ "crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"math/big"
	"strings"
)

// DecodedJWT is JWT that obtain oid
// all content base64 is decoded
type DecodedJWT struct {
	header    string
	payload   string
	signature string
}

// NewDecodedJwt return base64 decoded jwt
func NewDecodedJwt(data string) (*DecodedJWT, error) {
	decodes := []string{}
	for _, tkn := range strings.Split(data, ".") {
		str, err := base64.RawURLEncoding.DecodeString(tkn)
		if err != nil {
			return nil, err
		}
		decodes = append(decodes, string(str))
	}
	return &DecodedJWT{
		header:    decodes[0],
		payload:   decodes[1],
		signature: decodes[2],
	}, nil
}

// CalcHash calc hash for jwt signature
func (jwt *DecodedJWT) CalcHash() []byte {
	header := base64.RawURLEncoding.EncodeToString([]byte(jwt.header))
	payload := base64.RawURLEncoding.EncodeToString([]byte(jwt.payload))
	msg := header + "." + payload
	h := crypto.Hash.New(crypto.SHA256)
	h.Write([]byte(msg))
	return h.Sum(nil) //hashes
}

func getExponent(encoded string) (int, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		// err で判定するため、とりあえず０を返す
		return 0, err
	}
	if len(decoded) < 4 {
		// padding
		ndata := make([]byte, 4)
		copy(ndata[4-len(decoded):], decoded)
		decoded = ndata
	}
	return int(binary.BigEndian.Uint32(decoded)), nil
}

func getModulus(encoded string) (*big.Int, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	bi := new(big.Int)
	bi.SetBytes(decoded)
	return bi, nil
}

// NewPublicKey return `rsa.PublicKey`
func NewPublicKey(ee, em string) (*rsa.PublicKey, error) {
	e, err := getExponent(ee)
	if err != nil {
		return nil, err
	}

	bi, err := getModulus(em)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: bi,
		// E: 65537,
		E: e,
	}, nil
}
