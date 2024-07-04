package controllers

import (
	"encoding/json"
	"fmt"
	"go-backend/shared"
	"go-backend/types"
	"sort"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var playerListConnections = make(map[*websocket.Conn]bool)
var playerListLock = &sync.Mutex{}

func GetPlayers(c echo.Context) error {
	conn, err := shared.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func(conn *websocket.Conn) {
		playerListLock.Lock()
		delete(playerListConnections, conn)
		conn.Close()
		playerListLock.Unlock()
	}(conn)

	playerListLock.Lock()
	playerListConnections[conn] = true
	conn.WriteJSON(makePlayerList())
	playerListLock.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	return nil
}

func makePlayerList() [][]string {
	var players [][]string

	shared.Lock.RLock()
	defer shared.Lock.RUnlock()

	for playerName, token := range shared.PlayerNames {
		players = append(players, []string{playerName, token})
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i][0] < players[j][0]
	})

	return players
}

func BroadcastPlayerList() {
	for range shared.PlayerListChan {
		playerList := makePlayerList()
		playerListLock.Lock()
		for conn := range playerListConnections {
			go conn.WriteJSON(playerList)
		}
		playerListLock.Unlock()
	}
}

func UpdatePlayer(c echo.Context) error {
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

	fmt.Println(bodyJson["amount"])

	var pw string
	var ok bool
	if pw, ok = bodyJson["password"].(string); !ok {
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "No password provided",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}

	if pw != shared.Password {
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "Incorrect password",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}

	var amount string
	if amount, ok = bodyJson["amount"].(string); !ok {
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "No amount provided",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}

	var token string
	if token, ok = bodyJson["token"].(string); !ok {
		fmt.Println("No player selected")
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "No player selected",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}

	shared.Lock.Lock()
	var player types.Player
	if player, ok = shared.PlayerData[token]; !ok {
		fmt.Println("Player not found")
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "Player not found",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}

	amt, err := strconv.Atoi(amount)
	if err != nil {
		fmt.Println("Invalid amount")
	}
	player.Score += amt
	shared.PlayerData[token] = player
	shared.Lock.Unlock()

	return nil
}

func DeletePlayer(c echo.Context) error {
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

	var pw string
	var ok bool
	if pw, ok = bodyJson["password"].(string); !ok {
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "No password provided",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}

	if pw != shared.Password {
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "Incorrect password",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}

	var token string
	if token, ok = bodyJson["token"].(string); !ok {
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "No player selected",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}

	shared.Lock.Lock()
	var player types.Player
	if player, ok = shared.PlayerData[token]; !ok {
		enrichedJson, err := json.Marshal(map[string]string{
			"message": "Player not found",
			"success": "false",
		})
		if err != nil {
			return err
		}
		return c.JSONBlob(400, enrichedJson)
	}

	delete(shared.PlayerNames, player.Name)
	delete(shared.PlayerData, token)
	shared.Lock.Unlock()

	shared.PlayerListChan <- true

	return nil
}
