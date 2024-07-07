package shared

var (
	LeaderboardChan  = make(chan bool)
	BuzzedInChan     = make(chan bool)
	ResetBuzzersChan = make(chan bool)
	PlayerListChan   = make(chan bool)
	PlayersChan      = make(chan bool)
)
