package main

type ErrorResponse struct {
	Message string `json:"message"`
}

type ConfirmResponse struct {
	Token     string `json:"token"`
	ID        int64  `json:"id"`
	Userid    int64  `json:"userid"`
	Username  string `json:"username"`
	PublicKey string `json:"publickey"`
}

type SetUsernameResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type SetPublicKeyResponse struct {
	PublicKey string `json:"public_key"`
	Message   string `json:"message"`
}

type GetProfileResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

type SendMessageResponse struct {
	Message string `json:"message"`
}
