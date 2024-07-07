package controllers

import (
	"encoding/json"
	"go-backend/shared"
	"go-backend/types"
	"go-backend/util"
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
		return util.JsonParsingError(c)
	}

	var pw string
	var ok bool
	if pw, ok = bodyJson["password"].(string); !ok {
		return util.UserInputError(c, "No password provided")
	}

	if pw != shared.Password {
		return util.UserInputError(c, "Incorrect password")
	}

	var token string
	if token, ok = bodyJson["token"].(string); !ok || len(token) != 64 {
		return util.UserInputError(c, "No player selected")
	}

	var name string
	if name, ok = bodyJson["name"].(string); !ok || len(name) == 0 {
		return util.UserInputError(c, "No player name selected")
	}

	shared.Lock.Lock()
	defer shared.Lock.Unlock()
	var player types.Player
	if player, ok = shared.PlayerData[token]; !ok || player.Name != name {
		return util.UserInputError(c, "Player not found")
	}

	var amountStr string
	if amountStr, ok = bodyJson["amount"].(string); !ok {
		return util.UserInputError(c, "No amount provided")
	}
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return util.UserInputError(c, "Error parsing amount")
	}

	player.Score += int(amount)
	shared.PlayerData[token] = player

	shared.LeaderboardChan <- true

	enrichedJson, err := json.Marshal(map[string]string{
		"message": "Player score updated",
		"success": "true",
	})
	if err != nil {
		return err
	}

	return c.JSONBlob(200, enrichedJson)
}

func DeletePlayer(c echo.Context) error {
	bodyJson := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&bodyJson)
	if err != nil {
		return util.JsonParsingError(c)
	}

	var pw string
	var ok bool
	if pw, ok = bodyJson["password"].(string); !ok {
		return util.UserInputError(c, "No password provided")
	}

	if pw != shared.Password {
		return util.UserInputError(c, "Incorrect password")
	}

	var token string
	if token, ok = bodyJson["token"].(string); !ok {
		return util.UserInputError(c, "No player selected")
	}

	shared.Lock.Lock()
	var player types.Player
	if player, ok = shared.PlayerData[token]; !ok {
		return util.UserInputError(c, "Player not found")
	}

	delete(shared.PlayerNames, player.Name)
	delete(shared.PlayerData, token)
	shared.Lock.Unlock()

	shared.PlayerListChan <- true

	return nil
}
