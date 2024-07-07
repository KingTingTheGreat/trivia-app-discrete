package shared

import (
	"go-backend/types"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Lock     = &sync.RWMutex{}
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	PlayerNames    = make(map[string]string)
	PlayerData     = make(map[string]types.Player)
	QuestionNumber = 0
)
