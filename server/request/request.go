package request

type SetUsername struct {
	Username string `json:"username"`
}

type SetPublicKey struct {
	PublicKey string `json:"public_key"`
}

type SendMessage struct {
	Id         uint64 `json:"id"`
	Message    string `json:"message"`
	SessionKey string `json:"session_key"`
}

type ReplayMessage struct {
	MessageId uint64 `json:"message_id"`
	Message   string `json:"message"`
}
