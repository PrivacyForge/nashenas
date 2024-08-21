package request

type SetUsername struct {
	Username string `json:"username"`
}

type SetPublicKey struct {
	ReceivePublicKey string `json:"receive_public_key"`
	SendPublicKey    string `json:"send_public_key"`
}

type SendMessage struct {
	Id      uint64 `json:"id"`
	Message string `json:"message"`
}

type ReplyMessage struct {
	MessageId uint64 `json:"message_id"`
	Message   string `json:"message"`
}
