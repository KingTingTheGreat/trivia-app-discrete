package controllers

import (
	"errors"
	"fmt"
	"go-backend/shared"
	"go-backend/types"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func BuzzWs(c echo.Context) error {
	// upgrade the connection to a websocket connection
	conn, err := shared.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println("error upgrading")
		return err
	}
	// read the token from the connection
	_, tokenByte, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("error reading message")
		return err
	}
	token := string(tokenByte)

	// close the connection and remove websocket from player data
	defer func(conn *websocket.Conn) {
		shared.PlayerStore.PutPlayer(token, types.UpdatePlayer{Websocket: nil})
		conn.Close()
	}(conn)

	// check if the token is valid
	var player types.Player
	var ok bool
	if player, ok = shared.PlayerStore.GetPlayer(token); !ok {
		// return error
		return errors.New("invalid token")
	}

	// add the connection to the playerConnections map
	shared.PlayerStore.PutPlayer(token, types.UpdatePlayer{Websocket: conn})

	shared.LeaderboardChan <- true

	// read messages from the connection
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			break
		}
		message := string(p)
		fmt.Println(message)
		err = conn.WriteJSON(map[string]interface{}{
			"message": message,
		})
		if err != nil {
			break
		}
		// data race with store and player
		if player.BuzzedIn.IsZero() {
			player.BuzzedIn = time.Now()
			player.ButtonReady = false
			shared.PlayerStore.PutPlayer(token, types.UpdatePlayer{BuzzedIn: &player.BuzzedIn, ButtonReady: &player.ButtonReady})
			shared.BuzzedInChan <- true
		}
	}

	return nil
}
