package models

type User struct {
	ID     string `json:"id"`
	EMail  string `json:"email"`
	ApiKey string `json:"api_key"`
}
