package response

import "time"

type Error struct {
	Message string `json:"message"`
}

type Confirm struct {
	Token            string `json:"token"`
	ID               uint64 `json:"id"`
	Userid           uint64 `json:"userid"`
	Username         string `json:"username"`
	ReceivePublicKey string `json:"receive_public_key"`
	SendPublicKey    string `json:"send_public_key"`
}

type SetUsername struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type SetPublicKey struct {
	PublicKey string `json:"public_key"`
	Message   string `json:"message"`
}

type GetProfile struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

type SendMessage struct {
	Message   string `json:"message"`
	SessionID uint64 `json:"session_id"`
}

type GetMe struct {
	Username      string `json:"username"`
	Userid        uint64 `json:"userid"`
	PublicKey     string `json:"public_key"`
	PublicKeyHash string `json:"public_key_hash"`
}

type Quote struct {
	ID      uint64 `json:"id"`
	Content string `json:"content"`
}

type GetMessages struct {
	ID         uint64    `json:"id"`
	SessionID  uint64    `json:"session_id"`
	SessionKey string    `json:"session_key"`
	Content    string    `json:"content"`
	Time       time.Time `json:"time"`
	Owner      bool      `json:"owner"`
	Quote      *Quote    `json:"quote,omitempty"`
	CanReplay  bool      `json:"can_replay"`
}
