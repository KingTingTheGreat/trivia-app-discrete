package shared

import (
	"go-backend/types"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var (
	Lock     = &sync.RWMutex{}
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	PlayerNames      = make(map[string]string)
	PlayerData       = make(map[string]types.Player)
	QuestionNumber   = 0
	LeaderboardChan  = make(chan bool)
	BuzzedInChan     = make(chan bool)
	ResetBuzzersChan = make(chan bool)
	PlayerListChan   = make(chan bool)
	PlayersChan      = make(chan bool)
)

var Password string

func LoadPassword() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	Password = os.Getenv("PASSWORD")

	if Password == "" {
		panic("No password found in .env file")
	}
}
