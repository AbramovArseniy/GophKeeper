package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
)

func Encrypt(id string, MySecret []byte) (string, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return "", fmt.Errorf("unable to convert Id (NewUserSign func): %w", err)
	}

	uint32ID := uint32(intID)
	databyte := binary.BigEndian.AppendUint32(nil, uint32ID)
	hash := hmac.New(sha256.New, MySecret)
	hash.Write(databyte)

	sign := hash.Sum(nil)
	databyte = append(databyte, sign...)
	newsign := hex.EncodeToString(databyte)

	return newsign, nil
}

func Decrypt(s string, MySecret []byte) (string, bool, error) {
	decodedbyte, err := hex.DecodeString(s)
	if err != nil {
		return "", false, err
	}

	id := binary.BigEndian.Uint32(decodedbyte[:4])
	userID := strconv.Itoa(int(id))
	hash := hmac.New(sha256.New, MySecret)
	hash.Write(decodedbyte[:4])

	usersign := hash.Sum(nil)
	checkequal := hmac.Equal(usersign, decodedbyte[4:])

	return userID, checkequal, nil
}
