package controllers

import (
	"go-backend/shared"
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
}

func Reset(c echo.Context) error {
	ResetBuzzers()
	return c.JSONBlob(200, []byte(`{"success", "true"}`))
}
