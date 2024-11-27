package model

import "time"

// URL структура для JSON запросов и ответов.
type URL struct {
	Original string    `json:"original"`
	Short    string    `json:"short"`
	CreateTS time.Time `json:"create_ts"`
}
