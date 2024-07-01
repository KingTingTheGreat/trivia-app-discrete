package controllers

import (
	"fmt"
	"go-backend/shared"
	"go-backend/types"
	"sort"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var leaderboardConnections = make(map[*websocket.Conn]bool)
var leaderboardLock = &sync.Mutex{}

func Leaderboard(c echo.Context) error {
	conn, err := shared.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func(conn *websocket.Conn) {
		leaderboardLock.Lock()
		delete(leaderboardConnections, conn)
		conn.Close()
		leaderboardLock.Unlock()
	}(conn)

	leaderboardLock.Lock()
	leaderboardConnections[conn] = true
	leaderboardLock.Unlock()

	conn.WriteJSON(makeLeaderboard())

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	return nil
}

func makeLeaderboard() [][]string {
	var playerList []types.Player

	shared.Lock.RLock()
	defer shared.Lock.RUnlock()

	// get list of all players
	for _, player := range shared.PlayerData {
		playerList = append(playerList, player)
	}

	// sort the list by score, then last update, then name
	sort.Slice(playerList, func(i, j int) bool {
		if playerList[i].Score != playerList[j].Score {
			return playerList[i].Score > playerList[j].Score
		}
		if playerList[i].LastUpdate != playerList[j].LastUpdate {
			return playerList[i].LastUpdate.Before(playerList[j].LastUpdate)
		}
		return playerList[i].Name < playerList[j].Name
	})

	// create a list of player names and scores
	leaderboardList := make([][]string, 0)
	for _, player := range playerList {
		leaderboardList = append(leaderboardList, []string{player.Name, fmt.Sprintf("%d", player.Score)})
	}

	return leaderboardList
}

func BroadcastLeaderboard() {
	for range shared.LeaderboardChan {
		leaderboardLock.Lock()
		for conn := range leaderboardConnections {
			go conn.WriteJSON(makeLeaderboard())
		}
		leaderboardLock.Unlock()
	}
}
