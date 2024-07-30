package controllers

import (
	"encoding/json"
	"fmt"
	"go-backend/shared"
	"go-backend/types"
	"go-backend/util"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

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
		if player, ok := shared.PlayerStore.GetPlayer(token); ok && player.Name == name {
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
	}

	token, err = shared.PlayerStore.PostPlayer(types.Player{
		Name:             name,
		Score:            0,
		ButtonReady:      true,
		CorrectQuestions: make([]string, 0),
		LastUpdate:       time.Now(),
		BuzzedIn:         time.Time{},
		Websocket:        nil,
	})
	if err != nil {
		util.UserInputError(c, err.Error())
	}

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
