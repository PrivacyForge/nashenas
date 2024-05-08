package main

type SetUsernameRequest struct {
	Username string `json:"username"`
}

type SetPublicKeyRequest struct {
	PublicKey string `json:"public_key"`
}

type SendMessageRequest struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
}
