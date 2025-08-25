package client

import "fmt"

type QuestConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

func (q *QuestConfig) ToConnectionString() string {
	return fmt.Sprintf("http::addr=%s:%d;username=%s;password=%s;", q.Host, q.Port, q.Username, q.Password)
}
