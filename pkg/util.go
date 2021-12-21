package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

const (
	AES256Key = "AES256Key-32Characters1234567890"
)

var (
	errorKey = errors.New("key does not match")
)

func Encrypter(param string) (result string, err error) {

	key := []byte(AES256Key)
	plaintext := []byte(param)

	block, err := aes.NewCipher(key)
	if err != nil {
		return result, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(crand.Reader, nonce); err != nil {
		return result, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return result, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return fmt.Sprintf("%x", nonce) + `.` + fmt.Sprintf("%x", ciphertext), err
}

func Decrypter(param string) (result string, err error) {
	checkDot := strings.Contains(param, ".")
	if checkDot {
		params := strings.Split(param, ".")

		key := []byte(AES256Key)
		ciphertext, err := hex.DecodeString(params[1])
		if err != nil {
			return result, errorKey
		}
		nonce, err := hex.DecodeString(params[0])
		if err != nil {
			return result, errorKey
		}
		block, err := aes.NewCipher(key)
		if err != nil {
			// panic(err.Error())
			return result, errorKey
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			// panic(err.Error())
			return result, errorKey
		}

		plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			// panic(err.Error())
			return result, errorKey
		}

		result = fmt.Sprintf("%s", plaintext)
		return result, nil
	} else {
		return result, errorKey
	}
}

func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func TimePtr(t time.Time) *time.Time {
	return &t
}
