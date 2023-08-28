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

func NewHmacJwtMiddleware(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := ksuid.New().String()
		ctxWithRequestID := context.WithValue(c.Request.Context(), helper.ContextKeyRequestId, requestID)

		// Log request body
		logger := helper.GetLogger(ctxWithRequestID)
		if c.Request.Body != nil {
			body, _ := ioutil.ReadAll(c.Request.Body)
			logger.Info("Incoming request body ", string(body))
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
		}

		bearer := getBearerAuth(c.Request)
		if bearer == nil {
			c.Abort()
			response.WriteFailResponse(c, http.StatusUnauthorized, errors.New("missing bearer token"))
			return
		}
		claim, err := decodeHmacJwtData(secretKey, *bearer)
		if err != nil {
			c.Abort()
			response.WriteFailResponse(c, http.StatusUnauthorized, err)
			return
		}
		if err = claim.Validate(); err != nil {
			c.Abort()
			response.WriteFailResponse(c, http.StatusUnauthorized, err)
			return
		}
		c.Set(string(helper.ContextKeyJwtData), claim.User)
		c.Set(string(helper.ContextKeyTokenBearer), *bearer)
		c.Set(string(helper.ContextKeyRequestId), requestID)

		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()
		ctxWithJwt := context.WithValue(ctxWithRequestID, helper.ContextKeyJwtData, claim.User)
		ctxWithToken := context.WithValue(ctxWithJwt, helper.ContextKeyTokenBearer, *bearer)
		c.Request = c.Request.WithContext(ctxWithToken)

		c.Next()
	}
}

func GetJWTData(c *gin.Context) (model.User, error) {
	claim, err := c.Get(string(helper.ContextKeyJwtData))
	if !err {
		return model.User{}, model.NewNotFoundError()
	}
	return claim.(model.User), nil
}

func getBearerAuth(r *http.Request) *string {
	authHeader := r.Header.Get("Authorization")
	authForm := r.Form.Get("code")
	if authHeader == "" && authForm == "" {
		return nil
	}
	token := authForm
	if authHeader != "" {
		s := strings.SplitN(authHeader, " ", 2)
		if (len(s) != 2 || strings.ToLower(s[0]) != "bearer") && token == "" {
			return nil
		}
		//Use authorization header token only if token type is bearer else query string access token would be returned
		if len(s) > 0 && strings.ToLower(s[0]) == "bearer" {
			token = s[1]
		}
	}
	return &token
}

func decodeHmacJwtData(hmacSecret []byte, tokenStr string) (*JWTData, error) {
	var claim JWTData

	secretFn := func(token *jwt.Token) (interface{}, error) {
		if _, validSignMethod := token.Method.(*jwt.SigningMethodHMAC); !validSignMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	}

	token, err := jwt.ParseWithClaims(tokenStr, &claim, secretFn)
	if err != nil {
		return nil, err
	}

	if claim, ok := token.Claims.(*JWTData); ok && token.Valid {
		return claim, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GetToken(c context.Context) (token string, ok bool) {
	t := c.Value(helper.ContextKeyTokenBearer)
	if t != nil {
		token, ok = t.(string)
		return
	}
	t = c.Value(helper.ContextKeyTokenBearer.String())
	if t != nil {
		token, ok = t.(string)
		return
	}
	return
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
