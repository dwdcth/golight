package encrypt

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const (
	CHAR_SET               = "UTF-8"
	BASE_64_FORMAT         = "UrlSafeNoPadding"
	RSA_ALGORITHM_KEY_TYPE = "PKCS8"
	RSA_ALGORITHM_SIGN     = crypto.SHA256
)

// 返回生成的私钥、公钥对
func RSAKeyGen(bits int) (string, string, error) {
	//生成私钥
	privatekey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privatekey)
	block := &pem.Block{
		Type:  "RSA Private key",
		Bytes: derStream,
	}

	privateFile := bytes.NewBufferString("")
	err = pem.Encode(privateFile, block)
	if err != nil {
		return "", "", err
	}
	//利用私钥生成公钥
	publickey := &privatekey.PublicKey
	derpkix, err := x509.MarshalPKIXPublicKey(publickey)
	block = &pem.Block{
		Type:  "RSA Public key",
		Bytes: derpkix,
	}
	if err != nil {
		return "", "", err
	}
	publicFile := bytes.NewBufferString("")
	err = pem.Encode(publicFile, block)
	if err != nil {
		return "", "", err
	}
	return privateFile.String(), publicFile.String(), nil
}

// 公钥加密
func RsaPublicEncrypt(pubKey []byte, data []byte) ([]byte, error) {

	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("public key is bad")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil || pub == nil {
		return nil, errors.New("public key is bad")
	}
	pubRsa := pub.(*rsa.PublicKey)
	partLen := pubRsa.N.BitLen()/8 - 11
	chunks := splitArray(data, partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, pubRsa, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	}

	return buffer.Bytes(), nil
}

// 私钥解密
func RsaPrivateDecrypt(privateKey []byte, encrypted []byte) ([]byte, error) {

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key is bad")
	}

	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil || private == nil {
		return nil, errors.New("public key is bad")
	}
	privateRsa := private // private.(*rsa.PrivateKey)
	partLen := privateRsa.N.BitLen() / 8
	chunks := splitArray(encrypted, partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateRsa, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}

	return buffer.Bytes(), err
}

// 私钥数据加签
func RsaPrivateSign(privateKey []byte, data string) ([]byte, error) {

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key is bad")
	}

	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("public key is bad")
	}

	h := RSA_ALGORITHM_SIGN.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	sign, err := rsa.SignPKCS1v15(rand.Reader, private, RSA_ALGORITHM_SIGN, hashed)
	if err != nil {
		return nil, err
	}
	return sign, err
}

// 公钥数据验签
func RsaPubVerifySign(pubKey []byte, data string, sign string) error {

	block, _ := pem.Decode(pubKey)
	if block == nil {
		return errors.New("public key is bad")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil || pub == nil {
		return errors.New("public key is bad")
	}
	pubRsa := pub.(*rsa.PublicKey)
	h := RSA_ALGORITHM_SIGN.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	return rsa.VerifyPKCS1v15(pubRsa, RSA_ALGORITHM_SIGN, hashed, []byte(sign))
}

func splitArray(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}
