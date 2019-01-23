package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ossn/ossn-backend/helpers"
	"github.com/ossn/ossn-backend/middlewares"

	uuid "github.com/satori/go.uuid"

	"github.com/jinzhu/gorm"

	"github.com/ossn/ossn-backend/models"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	// OpenID Connect vars
	clientID     = os.Getenv("OPEN_ID_CLIENT_ID")
	clientSecret = os.Getenv("OPEN_ID_CLIENT_SECRET")
	providerURL  = os.Getenv("OPEN_ID_DOMAIN")
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
	ctx          context.Context
	provider     *oidc.Provider
	state        = helpers.RandStringBytesMaskImprSrc(15)
)

type Claims struct {
	Email       string `json:"email"`
	Verified    bool   `json:"email_verified"`
	Username    string `json:"nickname"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	AccessToken string `json:"access_token"`
}

func init() {
	ctx = context.Background()

	var err error
	provider, err = oidc.NewProvider(ctx, providerURL)
	if err != nil {
		log.Fatal(err)
	}
	verifier = provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  helpers.BackendURL + "oidc/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "openid"},
	}

}

func recoverFunc(w http.ResponseWriter, r *http.Request) {
	if rec := recover(); rec != nil {
		helpers.HandleError(w, r, http.StatusBadRequest, errors.New("Internal error"))
	}
}

func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	defer recoverFunc(w, r)
	if r.URL.Query().Get("state") != state {
		helpers.HandleError(w, r, http.StatusBadRequest, errors.New("State didn't match"))
		return
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		helpers.HandleError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		// handle missing token
		helpers.HandleError(w, r, http.StatusBadRequest, err)
		return
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		// handle error
		helpers.HandleError(w, r, http.StatusUnauthorized, err)
		return
	}

	// Extract custom claims
	claims := Claims{}
	if err := idToken.Claims(&claims); err != nil {
		helpers.HandleError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		helpers.HandleError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	user := &models.User{}
	userErr := models.DBSession.Where("oidc_id = ?", userInfo.Subject).First(user).Error
	if userErr != nil && userErr != gorm.ErrRecordNotFound {
		helpers.HandleError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Try to get github username if it doesn't exist
	if user.GithubURL == nil {
		c := make(chan string)
		go helpers.GetGithubURL(&userInfo.Subject, &claims.Username, c)
		defer func() {
			if user.GithubURL == nil {
				str := <-c
				if len(str) > 0 && user.ID > 0 {
					user.GithubURL = &str
					models.DBSession.Save(user)
				}
			}
		}()
	}

	user.Email = userInfo.Email
	user.ImageURL = &claims.Picture
	user.OIDCID = userInfo.Subject
	user.Name = claims.Name
	user.UserName = claims.Username

	models.DBSession.Save(user)

	token, err := middlewares.SignToken(user)
	if err != nil {
		helpers.HandleError(w, r, http.StatusInternalServerError, err)
		return
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		helpers.HandleError(w, r, http.StatusInternalServerError, err)
		return
	}
	cookie := uuid.String()

	session := &models.Session{
		UserID:      user.ID,
		Token:       token,
		Cookie:      cookie,
		AccessToken: oauth2Token.AccessToken,
	}

	err = models.DBSession.Create(session).Error
	if err != nil {
		helpers.HandleError(w, r, http.StatusInternalServerError, err)
		return
	}

	//TODO: Make this http only
	http.SetCookie(w, &http.Cookie{
		Expires: time.Now().Add(15 * time.Second),
		Name:    middlewares.AuthCookie,
		Value:   cookie,
	})

	http.Redirect(w, r, helpers.LoginURL+token, http.StatusTemporaryRedirect)
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
}
