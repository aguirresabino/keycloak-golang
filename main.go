package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"

	"golang.org/x/oauth2"
)

var (
	clientId     = "myclient"
	clientSecret = "176bbff7-7b04-4729-9cf1-38bc124b2b30"
	issuer       = "http://localhost:8080/auth/realms/Development"
	redirectUrl  = "http://localhost:8081/auth/callback"
)

func main() {

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectUrl,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := "123"

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/callback", (func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Query().Get("state") != state {
			http.Error(writer, "Invalid State", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(ctx, request.URL.Query().Get("code"))
		if err != nil {
			http.Error(writer, "Error exchange token", http.StatusInternalServerError)
			return
		}

		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(writer, "Failed to generate IDToken", http.StatusInternalServerError)
			return
		}

		result := struct {
			AccessToken *oauth2.Token
			IDToken     string
		}{
			token,
			idToken,
		}

		data, err := json.Marshal(result)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		writer.Write(data)
	}))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
