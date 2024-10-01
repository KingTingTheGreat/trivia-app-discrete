package shared

import (
	"errors"
	"go-backend/types"
	"math/rand"
	"sync"
)

// generates a random 64 character token
func generateToken() string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*_+-=<>?")
	var token string
	for i := 0; i < 64; i++ {
		token += string(chars[rand.Intn(len(chars))])
	}
	return token
}

type playerStore struct {
	mu          sync.RWMutex
	playerData  map[string]types.Player
	playerNames map[string]string
}

// returns existing player data
func (ps *playerStore) GetPlayer(token string) (types.Player, bool) {
	ps.mu.RLock()
	player, ok := ps.playerData[token]
	ps.mu.RUnlock()
	return player, ok
}

// creates a new player and returns the token
func (ps *playerStore) PostPlayer(player types.Player) (string, error) {
	if player.Name == "" {
		return "", errors.New("no player name provided")
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.playerNames[player.Name] != "" {
		return "", errors.New("a player with this name already exists")
	}

	token := generateToken()
	for _, ok := ps.playerData[token]; ok; {
		token = generateToken()
	}

	ps.playerData[token] = player
	ps.playerNames[player.Name] = token
	return token, nil
}

// updates an existing player
func (ps *playerStore) PutPlayer(token string, playerUpdates types.UpdatePlayer) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var player types.Player
	var ok bool
	// check this player exists
	if player, ok = ps.playerData[token]; !ok {
		return errors.New("player not found")
	}

	// update the player attributes
	if playerUpdates.ScoreDiff != nil {
		player.Score += *playerUpdates.ScoreDiff
	}
	if playerUpdates.ButtonReady != nil {
		player.ButtonReady = *playerUpdates.ButtonReady
	}
	if playerUpdates.LastUpdate != nil {
		player.LastUpdate = *playerUpdates.LastUpdate
	}
	if playerUpdates.BuzzedIn != nil {
		player.BuzzedIn = *playerUpdates.BuzzedIn
	}
	if playerUpdates.Websocket != nil {
		if player.Websocket != nil {
			player.Websocket.Close()
		}
		player.Websocket = playerUpdates.Websocket
	}
	ps.playerData[token] = player

	return nil
}

// deletes the player
func (ps *playerStore) DeletePlayer(token string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	var player types.Player
	var ok bool
	if player, ok = ps.playerData[token]; !ok {
		return errors.New("player not found")
	}
	player.Websocket.Close()
	name := ps.playerData[token].Name
	delete(ps.playerData, token)
	delete(ps.playerNames, name)
	return nil
}

// returns a list of all players
func (ps *playerStore) AllPlayers() []types.Player {
	var allPlayers []types.Player
	ps.mu.RLock()
	for _, player := range ps.playerData {
		allPlayers = append(allPlayers, player)
	}
	ps.mu.RUnlock()
	return allPlayers
}

type TokenPlayer struct {
	Token  string
	Player types.Player
}

// returns a list of all tokens and their corresponding players
func (ps *playerStore) AllTokenPlayers() []TokenPlayer {
	var allPlayers []TokenPlayer
	ps.mu.RLock()
	for token, player := range ps.playerData {
		allPlayers = append(allPlayers, TokenPlayer{token, player})
	}
	ps.mu.RUnlock()
	return allPlayers
}

type PlayerNameToken struct {
	Name string 
	Token string
}

// returns a list of all player names and their corresponding tokens
func (ps *playerStore) AllNamesTokens() []PlayerNameToken {
	var allNamesTokens []PlayerNameToken
	ps.mu.RLock()
	for name, token := range ps.playerNames {
		allNamesTokens = append(allNamesTokens, PlayerNameToken{
			Name:name, 
			Token: token,
		})
	}
	ps.mu.RUnlock()
	return allNamesTokens
}

// returns the token for a given player name
func (ps *playerStore) NameToToken(name string) (string, bool) {
	ps.mu.RLock()
	token, ok := ps.playerNames[name]
	ps.mu.RUnlock()
	return token, ok
}

// returns a boolean indicating if the token and name match
func (ps *playerStore) VerifyTokenName(token, name string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	storedToken, ok := ps.playerNames[name]
	return ok && storedToken == token
}

var PlayerStore playerStore = playerStore{
	mu:          sync.RWMutex{},
	playerData:  make(map[string]types.Player),
	playerNames: make(map[string]string),
}
