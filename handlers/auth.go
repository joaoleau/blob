package handlers

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func InitAuth(config *oauth2.Config, state string) gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectURL := config.AuthCodeURL(state)
		c.Redirect(http.StatusFound, redirectURL)
	}
}

func CallbackAuth(config *oauth2.Config, ctx context.Context, state string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Query("state") != state {
			c.JSON(http.StatusBadRequest, gin.H{"error": "State doesn't match"})
			return
		}

		oauth2Token, err := config.Exchange(ctx, c.Query("code"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Problema ao trocar Token", "details": err.Error()})
			return
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Problema ao pegar ID Token"})
			return
		}

		res := struct {
			OAuth2Token *oauth2.Token `json:"oauth2_token"`
			IDToken     string        `json:"id_token"`
		}{
			OAuth2Token: oauth2Token,
			IDToken:     rawIDToken,
		}

		c.JSON(http.StatusOK, res)
	}
}