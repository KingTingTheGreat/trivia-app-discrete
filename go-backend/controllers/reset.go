package controllers

import (
	"encoding/json"
	"fmt"
	"go-backend/shared"
	"go-backend/util"
	"time"

	"github.com/labstack/echo/v4"
)

func ResetBuzzersLoop() {
	for range shared.ResetBuzzersChan {
		ResetBuzzers()
	}
}

func ResetBuzzers() {
	shared.Lock.Lock()
	for key, player := range shared.PlayerData {
		player.ButtonReady = true
		if player.Websocket == nil {
			continue
		}
		err := player.Websocket.WriteJSON(map[string]interface{}{
			"buttonReady": "true",
		})
		if err != nil {
			player.Websocket.Close()
			player.Websocket = nil
		}
		player.BuzzedIn = time.Time{}
		shared.PlayerData[key] = player
	}
	shared.Lock.Unlock()

	shared.BuzzedInChan <- true
}

func Reset(c echo.Context) error {
	bodyJson := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&bodyJson)
	if err != nil {
		fmt.Println("error parsing")
		return util.JsonParsingError(c)
	}

	var pw string
	var ok bool
	if pw, ok = bodyJson["password"].(string); !ok || len(pw) == 0 {
		fmt.Println("no password")
		return util.UserInputError(c, "No password provided")
	}

	if pw != shared.Password {
		fmt.Println("incorrect password")
		return util.UserInputError(c, "Incorrect password")
	}

	ResetBuzzers()

	enrichedJson, err := json.Marshal(map[string]string{
		"success": "true",
	})
	if err != nil {
		return err
	}

	return c.JSONBlob(200, enrichedJson)
}
