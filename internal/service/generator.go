package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"url-shortener/internal/helper"
)

var charset string = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func generate() (string, error) {

	slice := make([]byte, 6)
	for i := range slice {
		newBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			fmt.Println("(-_-) как такое возможно...")
			return "", helper.ErrInternalServer
		}
		slice[i] = charset[newBig.Int64()]
	}

	return string(slice), nil

}

func GenerateAndCheck(ctx context.Context, s *URLService) (string, error) {

	for range 5 {

		code, err := generate()
		if err != nil {
			return "", err
		}

		exist, err := s.repo.CheckCode(ctx, code)
		if err != nil {
			return "", err
		}

		if !exist {
			return code, nil
		}

	}

	return "", errors.New("Layer: generator.go. CheckAndGenerate error (end cycle with no result)")

}
