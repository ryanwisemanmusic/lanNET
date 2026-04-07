package proto

import (
	"encoding/json"
)

const MaxPlayers = 4

type MessageType string

const (
	TypeJoin        MessageType = "join"
	TypeJoinAccept  MessageType = "join_accept"
	TypeJoinReject  MessageType = "join_reject"
	TypePlayerState MessageType = "player_state"
	TypeChat        MessageType = "chat"
	TypeLeave       MessageType = "leave"
)

type Message struct {
	Type    MessageType `json:"type"`
	FromID  int         `json:"from_id,omitempty"`
	ToID    int         `json:"to_id,omitempty"`
	Payload []byte      `json:"payload"`
}

type JoinPayload struct {
	PlayerName string `json:"player_name"`
}

type JoinAcceptPayload struct {
	AssignedID int    `json:"assigned_id"`
	PlayerIDs  [4]int `json:"player_ids"`
}

type PlayerStatePaylooad struct {
	X     float32 `json:"x,omitempty"`
	Y     float32 `json:"y,omitempty"`
	HP    int     `json:"hpm,omitempty"`
	Score int     `json:"score,omitempty"`
}

func EncodeMessage(msg *Message) ([]byte, error) {
	return json.Marshal(msg)
}

func DecodeMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
