package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func Encrypt(data, MySecret []byte) ([]byte, error) {
	block, err := aes.NewCipher(MySecret)
	if err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, make([]byte, block.BlockSize()))
	cipherText := make([]byte, len(data))
	cfb.XORKeyStream(cipherText, data)
	return cipherText, nil
}

func Decrypt(data, MySecret []byte) ([]byte, error) {
	block, err := aes.NewCipher(MySecret)
	if err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBDecrypter(block, make([]byte, block.BlockSize()))
	plainText := make([]byte, len(data))
	cfb.XORKeyStream(plainText, data)
	return plainText, nil
}
