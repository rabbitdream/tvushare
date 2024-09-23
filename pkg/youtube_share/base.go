package youtube_share

import (
	"log"
	"tvushare/model"

	"golang.org/x/oauth2"
)

func exchangeTokenV2(config *oauth2.Config, code string) (*model.OauthToken, error) {
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token %v", err)
	}
	token := &model.OauthToken{
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		ExpiresIn:    int(tok.Extra("expires_in").(float64)),
		Scope:        tok.Extra("scope").(string),
	}
	return token, nil
}
