package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
)

func RsaGenerateKey(size int) (pubkey string, prikey string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return
	}
	// 生成私钥
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prikey = string(pem.EncodeToMemory(block))
	// 生成公钥
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = string(pem.EncodeToMemory(block))
	return
}

// 格式化秘钥，去除头尾和换行
func RsaFormatKeyToString(key string) (fkey string) {
	fkey = strings.Replace(key, " ", "", -1)
	fkey = strings.Replace(fkey, "\n", "", -1)
	fkey = strings.Replace(fkey, "-----BEGINPUBLICKEY-----", "", 1)
	fkey = strings.Replace(fkey, "-----ENDPUBLICKEY-----", "", 1)
	fkey = strings.Replace(fkey, "-----BEGINPRIVATEKEY-----", "", 1)
	fkey = strings.Replace(fkey, "-----ENDPRIVATEKEY-----", "", 1)
	fkey = strings.Replace(fkey, "-----BEGINRSAPRIVATEKEY-----", "", 1)
	fkey = strings.Replace(fkey, "-----ENDRSAPRIVATEKEY-----", "", 1)
	return
}

// 格式化秘钥，添加头尾
func RsaFormatKeyToPem(key string, isPub bool) (pkey string) {
	pkey = WrapString(strings.TrimSpace(key), 64)
	if isPub {
		pkey = "-----BEGIN PUBLIC KEY-----\n" + pkey + "\n-----END PUBLIC KEY-----\n"
	} else {
		pkey = "-----BEGIN RSA PRIVATE KEY-----\n" + pkey + "\n-----END RSA PRIVATE KEY-----\n"
	}
	return
}

// 加密
func RsaEncrypt(origData string, pubkey string) string {
	origByte := []byte(origData)
	publicKey := []byte(pubkey)
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return ""
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return ""
	}
	pub := pubInterface.(*rsa.PublicKey)
	strEncryptByte, err := rsa.EncryptPKCS1v15(rand.Reader, pub, origByte)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(strEncryptByte)
}

// 解密
func RsaDecrypt(str string, prikey string) string {
	ciphertext, _ := base64.StdEncoding.DecodeString(str)
	privateKey := []byte(prikey)
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return ""
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}
	//strDecryptByte, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, priv, ciphertext, nil)
	strDecryptByte, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return ""
	}
	return string(strDecryptByte)
}

func RsaParsePublicKey(str string) *rsa.PublicKey {
	str = strings.TrimSpace(str)
	if !strings.HasPrefix(str, "-----BEGIN PUBLIC KEY-----") {
		str = RsaFormatKeyToPem(str, true)
	}
	publicKey := []byte(str)
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil
	}
	pub := pubInterface.(*rsa.PublicKey)
	return pub
}

func RsaParsePrivateKey(str string) *rsa.PrivateKey {
	str = strings.TrimSpace(str)
	if !strings.HasPrefix(str, "-----BEGIN RSA PRIVATE KEY-----") {
		str = RsaFormatKeyToPem(str, false)
	}
	privateKey := []byte(str)
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}
	return priv
}
