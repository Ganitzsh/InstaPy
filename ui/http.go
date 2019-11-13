package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func jsonError(c *gin.Context, status int, err error) {
	c.JSON(status, map[string]string{
		"error": err.Error(),
	})
}

func validateToken(c *gin.Context) {
	bearer := strings.Split(c.GetHeader("Authorization"), " ")
	if len(bearer) < 2 {
		jsonError(c, http.StatusUnauthorized, errors.New("Invalid Token"))
		return
	}
	bearerToken := bearer[1]
	if err := verifyToken(&token{accessToken: bearerToken}); err != nil {
		jsonError(c, http.StatusUnauthorized, errors.New("Invalid Token"))
		return
	}
	c.Next()
}

type authenticatePayload struct {
	username string
	password string
}

func authenticate(c *gin.Context) {
	req := authenticatePayload{}

	if err := c.BindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}

	t, err := authenticateUser(req.username, req.password)
	if err != nil {
		jsonError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, t)
}
