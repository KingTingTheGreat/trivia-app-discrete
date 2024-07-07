package controllers

import (
	"encoding/json"
	"go-backend/shared"
	"go-backend/types"
	"go-backend/util"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// checks if the token provided matches the name provided
func PostVerify(c echo.Context) error {
	// ready request body json
	bodyJson := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&bodyJson)
	if err != nil {
		return util.JsonParsingError(c)
	}
	// check if the request body contains the correct key
	var token string
	var ok bool
	if token, ok = bodyJson["token"].(string); !ok {
		return c.JSON(400, "No token provided")
	}
	var name string
	if name, ok = bodyJson["name"].(string); !ok {
		return c.JSON(400, "No name provided")
	}
	name = strings.TrimSpace(name)

	var player types.Player
	shared.Lock.RLock()
	defer shared.Lock.RUnlock()
	if player, ok = shared.PlayerData[token]; !ok {
		return c.JSON(400, "No player with this token exists")
	}

	if player.Name != name {
		return c.JSON(400, "The name provided does not match the token")
	}

	enrichedJson, err := json.Marshal(map[string]string{
		"success":     "true",
		"buttonReady": strconv.FormatBool(player.ButtonReady),
	})
	if err != nil {
		return err
	}

	return c.JSONBlob(200, enrichedJson)
}
