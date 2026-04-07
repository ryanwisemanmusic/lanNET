package session

import (
	"fmt"
	"lanNET/internal/proto"
	"net"
	"sync"
	"time"
)

type Player struct {
	ID       int
	Name     string
	Conn     net.Conn
	LastSeen time.Time
	isReady  bool
}

type SessionManager struct {
	mu          sync.RWMutex
	players     [proto.MaxPlayers]*Player
	playerCount int
	lobbyReady  bool
}

func NewMessionManager() *SessionManager {
	return &SessionManager{}
}

func (sm *SessionManager) AddPlayer(conn net.Conn, name string) (int, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.playerCount >= proto.MaxPlayers {
		return -1, fmt.Errorf("session full: max %d players", proto.MaxPlayers)
	}
	id := -1
	for i := 0; i < proto.MaxPlayers; i++ {
		if sm.players[i] == nil {
			id = i
			break
		}
	}

	if id == -1 {
		return -1, fmt.Errorf("no available slot")
	}

	sm.players[id] = &Player{
		ID:       id,
		Name:     name,
		Conn:     conn,
		LastSeen: time.Now(),
	}
	sm.playerCount++

	return id, nil
}

func (sm *SessionManager) RemovePlayer(id int) {
	sm.mu.Lock()
	defer sm.mu.Lock()

	if sm.players[id] != nil {
		sm.players[id] = nil
		sm.playerCount--
	}
}

func (sm *SessionManager) GetPlayer(id int) *Player {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.players[id]
}

func (sm *SessionManager) GetAllPlayers() []*Player {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	players := make([]*Player, 0, sm.playerCount)
	for _, p := range sm.players {
		if p != nil {
			players = append(players, p)
		}
	}
	return players
}

func (sm *SessionManager) IsFull() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.playerCount >= proto.MaxPlayers
}

func (sm *SessionManager) GetPlayerIDs() [4]int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var ids [4]int
	for i := 0; i < proto.MaxPlayers; i++ {
		if sm.players[i] != nil {
			ids[i] = sm.players[i].ID
		} else {
			ids[i] = -1
		}
	}
	return ids
}

func (sm *SessionManager) BroadcastToAll(msg []byte, excludeID int) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for _, p := range sm.players {
		if p != nil && p.ID != excludeID {
			p.Conn.Write(msg)
		}
	}
}
