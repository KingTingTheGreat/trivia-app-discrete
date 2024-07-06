package controllers

import (
	"encoding/json"
	"go-backend/shared"
	"go-backend/types"
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
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "Error parsing request body. Please try again",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}
	// check if the request body contains the correct key
	var name string
	var ok bool
	if name, ok = bodyJson["name"].(string); !ok {
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "No name provided",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}
	name = strings.TrimSpace(name)

	shared.Lock.Lock()
	if _, ok = shared.PlayerNames[name]; ok {
		shared.Lock.Unlock()
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "A player with this name already exists",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}

	token := generateToken()
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
