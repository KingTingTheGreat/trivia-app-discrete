package controllers

import (
	"encoding/json"
	"fmt"
	"go-backend/shared"
	"go-backend/types"
	"go-backend/util"
	"math/rand"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func generateToken() string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*_+-=<>?")
	var token string
	for i := 0; i < 64; i++ {
		token += string(chars[rand.Intn(len(chars))])
	}
	return token
}

// returns a new token for a player
func PostToken(c echo.Context) error {
	// read request body json
	bodyJson := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&bodyJson)
	if err != nil {
		return util.JsonParsingError(c)
	}
	// check if the request body contains the correct key
	var name string
	var ok bool
	if name, ok = bodyJson["name"].(string); !ok {
		return util.UserInputError(c, "No name provided")
	}
	name = strings.TrimSpace(name)

	fmt.Println(bodyJson)
	var token string
	if token, ok = bodyJson["token"].(string); ok && len(token) == 64 {
		// check if this player and token already exist
		shared.Lock.RLock()
		if player, ok := shared.PlayerData[token]; ok && player.Name == name && shared.PlayerNames[name] == token {
			shared.Lock.RUnlock()
			enrichedJson, err := json.Marshal(map[string]string{
				"message": "Successfully restored player",
				"success": "true",
				"token":   token,
				"name":    name,
			})
			if err != nil {
				return err
			}
			return c.JSONBlob(200, enrichedJson)
		}
		shared.Lock.RUnlock()
	}

	shared.Lock.Lock()
	if _, ok = shared.PlayerNames[name]; ok {
		shared.Lock.Unlock()
		return util.UserInputError(c, "A player with this name already exists")
	}

	if len(token) != 64 {
		token = generateToken()
	}
	for _, ok = shared.PlayerData[token]; ok; {
		token = generateToken()
	}

	shared.PlayerNames[name] = token
	shared.PlayerData[token] = types.Player{
		Name:             name,
		Score:            0,
		ButtonReady:      false,
		CorrectQuestions: make([]string, 0),
		LastUpdate:       time.Now(),
		BuzzedIn:         time.Time{},
		Websocket:        nil,
	}
	shared.Lock.Unlock()

	shared.PlayerListChan <- true
	shared.LeaderboardChan <- true

	enrichedJson, err := json.Marshal(map[string]string{
		"message": "Token generated successfully",
		"success": "true",
		"token":   token,
		"name":    name,
	})
	if err != nil {
		return err
	}
	return c.JSONBlob(200, enrichedJson)
}
