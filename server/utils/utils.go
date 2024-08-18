package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	_stringProps = map[string]bool{
		"start_param": true,
	}
)

type Chat struct {
	ID int64 `json:"id"`
	Type string `json:"type"`
	Title string `json:"title"`
	PhotoURL string `json:"photo_url"`
	Username string `json:"username"`
}

type User struct {
	AddedToAttachmentMenu bool `json:"added_to_attachment_menu"`
	AllowsWriteToPm bool `json:"allows_write_to_pm"`
	FirstName string `json:"first_name"`
	ID int64 `json:"id"`
	IsBot bool `json:"is_bot"`
	IsPremium bool `json:"is_premium"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	LanguageCode string `json:"language_code"`
	PhotoURL string `json:"photo_url"`
}

type InitData struct {
	AuthDateRaw int `json:"auth_date"`
	CanSendAfterRaw int `json:"can_send_after"`
	Chat Chat `json:"chat"`
	ChatType string `json:"chat_type"`
	ChatInstance int64 `json:"chat_instance"`
	Hash string `json:"hash"`
	QueryID string `json:"query_id"`
	Receiver User `json:"receiver"`
	StartParam string `json:"start_param"`
	User User `json:"user"`
}

func Parse(initData string) (InitData, error) {
	q, err := url.ParseQuery(initData)
	if err != nil {
		return InitData{}, errors.New("format error")
	}

	pairs := make([]string, 0, len(q))
	for k, v := range q {

		val := v[0]
		valFormat := "%q:%q"
		if isString := _stringProps[k]; !isString && json.Valid([]byte(val)) {
			valFormat = "%q:%s"
		}

		pairs = append(pairs, fmt.Sprintf(valFormat, k, val))
	}

	var d InitData
	jStr := fmt.Sprintf("{%s}", strings.Join(pairs, ","))
	if err := json.Unmarshal([]byte(jStr), &d); err != nil {
		return InitData{}, errors.New("format error")
	}
	return d, nil
}
