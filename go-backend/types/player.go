package types

import (
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	Name             string
	Score            int
	ButtonReady      bool
	CorrectQuestions []string
	LastUpdate       time.Time
	BuzzedIn         time.Time
	Websocket        *websocket.Conn
}

type UpdatePlayer struct {
	Name        *string
	ScoreDiff   *int
	ButtonReady *bool
	LastUpdate  *time.Time
	BuzzedIn    *time.Time
	Websocket   *websocket.Conn
}

type BuzzedInPlayer struct {
	Name string
	Time string
}

type LeaderboardPlayer struct {
	Name  string
	Score int
}
