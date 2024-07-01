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
