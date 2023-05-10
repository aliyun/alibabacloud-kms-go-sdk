package utils

import (
	"bytes"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

func DecoderBytesToString(data []byte, encoding string) (string, error) {
	encoder, err := ianaindex.IANA.Encoding(encoding)
	if err != nil {
		return "", err
	}
	reader := transform.NewReader(bytes.NewReader(data), encoder.NewDecoder())
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func EncoderStringToBytes(plaintext string, encoding string) ([]byte, error) {
	encoder, err := ianaindex.IANA.Encoding(encoding)
	if err != nil {
		return nil, err
	}
	reader := transform.NewReader(strings.NewReader(plaintext), encoder.NewEncoder())
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return result, nil

}
