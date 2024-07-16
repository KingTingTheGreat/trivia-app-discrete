package controllers

import (
	"fmt"
	"go-backend/shared"
	"sort"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var buzzedInConnections = make(map[*websocket.Conn]bool)
var buzzedInLock = &sync.Mutex{}

func BuzzedIn(c echo.Context) error {
	conn, err := shared.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func(conn *websocket.Conn) {
		buzzedInLock.Lock()
		delete(buzzedInConnections, conn)
		conn.Close()
		buzzedInLock.Unlock()
	}(conn)

	buzzedInLock.Lock()
	buzzedInConnections[conn] = true
	conn.WriteJSON(makeBuzzedIn())
	buzzedInLock.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	return nil
}

func makeBuzzedIn() [][]string {
	playerList := shared.PlayerStore.AllPlayers()

	// sort the list by buzz in time, then score, then name
	sort.Slice(playerList, func(i, j int) bool {
		if playerList[i].BuzzedIn != playerList[j].BuzzedIn {
			return playerList[i].BuzzedIn.Before(playerList[j].BuzzedIn)
		}
		if playerList[i].Score != playerList[j].Score {
			return playerList[i].Score > playerList[j].Score
		}
		return playerList[i].Name < playerList[j].Name
	})

	// create a list of player names and buzz in times
	buzzedInList := make([][]string, 0)
	for _, player := range playerList {
		// filter out players who haven't buzzed in
		if player.BuzzedIn.IsZero() {
			continue
		}
		buzzedInList = append(buzzedInList, []string{player.Name, player.BuzzedIn.Format("03:04:05.000 PM")})
	}

	return buzzedInList
}

func BroadcastBuzzedIn() {
	for range shared.BuzzedInChan {
		buzzedIn := makeBuzzedIn()
		fmt.Println(("broadcasting buzzed in"))
		buzzedInLock.Lock()
		for conn := range buzzedInConnections {
			go conn.WriteJSON(buzzedIn)
		}
		buzzedInLock.Unlock()
	}
}
