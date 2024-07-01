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
	PlayerNames      = make(map[string]bool)
	PlayerData       = make(map[string]types.Player)
	QuestionNumber   = 0
	LeaderboardChan  = make(chan bool)
	BuzzedInChan     = make(chan bool)
	ResetBuzzersChan = make(chan bool)
	PlayersChan      = make(chan bool)
)
