package handler

import (
	"log"
	"net/http"

	"github.com/Mobrick/gophermart/internal/config"
	"github.com/Mobrick/gophermart/internal/database"
	"github.com/Mobrick/gophermart/internal/userauth"
)

type HandlerEnv struct {
	ConfigStruct *config.Config
	Storage      database.Storage
}

func GetUserIDFromRequest(req *http.Request) (string, bool) {
	cookie, err := req.Cookie("auth_token")
	if err != nil {
		log.Printf("no cookie found. " + err.Error())
		return "", false
	}	

	token := cookie.Value
	userID, ok := userauth.GetUserID(token)
	if !ok {
		log.Printf("invalid token")
		return "", false
	}
	return userID, true
}
