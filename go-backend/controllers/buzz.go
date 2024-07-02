package controllers

import (
	"errors"
	"fmt"
	"go-backend/shared"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func BuzzWs(c echo.Context) error {
	// upgrade the connection to a websocket connection
	conn, err := shared.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	// read the token from the connection
	_, tokenByte, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	token := string(tokenByte)

	// close the connection and remove websocket from player data
	defer func(conn *websocket.Conn) {
		shared.Lock.Lock()
		player := shared.PlayerData[token]
		player.Websocket = nil
		shared.PlayerData[token] = player
		conn.Close()
		shared.Lock.Unlock()
	}(conn)

	// check if the token is valid
	if _, ok := shared.PlayerData[token]; !ok {
		// return error
		return errors.New("invalid token")
	}

	// add the connection to the playerConnections map
	shared.Lock.Lock()
	player := shared.PlayerData[token]
	player.Websocket = conn
	shared.PlayerData[token] = player
	shared.Lock.Unlock()

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
		if player.BuzzedIn.IsZero() {
			shared.Lock.Lock()
			player.BuzzedIn = time.Now()
			shared.PlayerData[token] = player
			shared.Lock.Unlock()
			shared.BuzzedInChan <- true
		}
	}

	return nil
}
