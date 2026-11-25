package user

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Mail      string    `json:"mail"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
