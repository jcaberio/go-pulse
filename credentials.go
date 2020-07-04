package pulse

import (
	"encoding/base64"
	"encoding/json"
)

type credentials struct {
	username string
	password string
}

func newCredentials(username, password string) *credentials {
	return &credentials{
		username: base64.StdEncoding.EncodeToString([]byte(username)),
		password: base64.StdEncoding.EncodeToString([]byte(password)),
	}
}

func (request *credentials) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: request.username,
		Password: request.password,
	})
}
