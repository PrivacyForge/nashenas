package response

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
	Message string `json:"message"`
}

type GetMe struct {
	ID               uint64 `json:"id"`
	Username         string `json:"username"`
	Userid           uint64 `json:"userid"`
	ReceivePublicKey string `json:"receive_public_key"`
	SendPublicKey    string `json:"send_public_key"`
}
