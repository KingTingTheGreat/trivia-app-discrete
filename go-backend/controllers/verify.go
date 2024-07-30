package controllers

import (
	"encoding/json"
	"fmt"
	"go-backend/shared"
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

	if !shared.PlayerStore.VerifyTokenName(token, name) {
		return c.JSON(400, "The token and name provided do not match")
	}

	player, _ := shared.PlayerStore.GetPlayer(token)

	enrichedJson, err := json.Marshal(map[string]string{
		"success":     "true",
		"buttonReady": strconv.FormatBool(player.ButtonReady),
	})
	if err != nil {
		return err
	}

	fmt.Println(player.ButtonReady)

	return c.JSONBlob(200, enrichedJson)
}
