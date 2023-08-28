package middleware

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"tempo/controller/response"
	"tempo/helper"
	"tempo/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type JWTData struct {
	jwt.StandardClaims
	model.User
}

func GenerateJwt(user model.User, secretKey string) (*string, error) {
	stdClaims := jwt.StandardClaims{
		Id:        ksuid.New().String(),
		ExpiresAt: time.Now().Add(300 * time.Second).Unix(),
		Subject:   *user.Id,
	}

	// generate token
	accessClaims := JWTData{
		StandardClaims: stdClaims,
		User:           user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	var byteSecret = []byte(secretKey)
	accessToken, err := token.SignedString(byteSecret)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}
