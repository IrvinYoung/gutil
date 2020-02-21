package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"hash"
)

func LoadRSAPrivateKey(d []byte) (priv *rsa.PrivateKey, err error) {
	p, _ := pem.Decode(d)
	if p == nil {
		err = errors.New("load private key failed")
		return
	}
	priv, err = x509.ParsePKCS1PrivateKey(p.Bytes)
	return
}

func LoadRSAPublicKey(d []byte) (pub *rsa.PublicKey, err error) {
	p, _ := pem.Decode(d)
	if p == nil {
		err = errors.New("load public key failed")
		return
	}
	pi, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		return
	}
	pub = pi.(*rsa.PublicKey)
	return
}

func MaxRSAEncryptLengthByOAEP(hash hash.Hash, pub *rsa.PublicKey) int {
	if pub == nil {
		return 0
	}
	k := (pub.N.BitLen() + 7) / 8
	return k - 2*hash.Size() - 2
}

func MaxRSADecryptLength(priv *rsa.PrivateKey) int {
	if priv == nil {
		return 0
	}
	return (priv.PublicKey.N.BitLen() + 7) / 8
}

func RSAEncrypt(pub *rsa.PublicKey, origData []byte) (cipherData []byte, err error) {
	if pub == nil || len(origData) == 0 {
		err = errors.New("params error")
		return
	}

	var (
		hash               = sha256.New()
		buf                []byte
		maxLen, begin, end int
		totalLen           = len(origData)
		goon               = true
	)
	if maxLen = MaxRSAEncryptLengthByOAEP(hash, pub); maxLen == 0 {
		err = rsa.ErrMessageTooLong
		return
	}
	for i := 0; goon; i++ {
		begin, end = i*maxLen, (i+1)*maxLen
		if end > totalLen {
			buf = origData[begin:]
			goon = false
		} else {
			buf = origData[begin:end]
		}
		if buf, err = rsa.EncryptOAEP(hash, rand.Reader, pub, buf, nil); err != nil {
			return
		}
		cipherData = append(cipherData, buf...)
	}
	return
}

func RSADecrypt(priv *rsa.PrivateKey, cipherData []byte) (origData []byte, err error) {
	if priv == nil || len(cipherData) == 0 {
		err = errors.New("params error")
		return
	}

	var (
		hash               = sha256.New()
		buf                []byte
		maxLen, begin, end int
		totalLen           = len(cipherData)
		goon               = true
	)

	if maxLen = MaxRSADecryptLength(priv); maxLen == 0 {
		err = rsa.ErrMessageTooLong
		return
	}
	for i := 0; goon; i++ {
		begin, end = i*maxLen, (i+1)*maxLen
		if end >= totalLen {
			buf = cipherData[begin:]
			goon = false
		} else {
			buf = cipherData[begin:end]
		}
		if buf, err = rsa.DecryptOAEP(hash, rand.Reader, priv, buf, nil); err != nil {
			return
		}
		origData = append(origData, buf...)
	}
	return
}
