package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

func Generate(text, secretKey string) (string, error) {
	if text == "" || secretKey == "" {
		return "", errors.New("text and secretKey cannot be empty")
	}

	encodedText := base64.StdEncoding.EncodeToString([]byte(text))

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(text))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return encodedText + "." + signature, nil
}

func Verify(signedText string, secretKey string) (bool, string, error) {
	parts := strings.Split(signedText, ".")
	if len(parts) != 2 {
		return false, "", errors.New("invalid signed text format")
	}

	textBytes, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, "", err
	}
	originalText := string(textBytes)

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(originalText))
	expectedSignature := h.Sum(nil)

	inputSignature, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, "", err
	}

	if !hmac.Equal(inputSignature, expectedSignature) {
		return false, "", nil
	}

	return true, originalText, nil
}
