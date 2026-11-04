package post

import (
	"time"

	"github.com/CALEXCO/minisocial/internal/user"
)

type Post struct {
	user.User `json:"user"`
	Msg       string    `json:"msg"`
	Date      time.Time `json:"date"`
}
