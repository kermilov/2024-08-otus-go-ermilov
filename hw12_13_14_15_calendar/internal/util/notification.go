package util

import "time"

type Notification struct {
	// ID события;
	ID string `json:"id"`
	// Заголовок события;
	Title string `json:"title"`
	// Дата события;
	DateTime time.Time `json:"dateTime"`
	// Пользователь, которому отправлять.
	UserID int64 `json:"userId"`
}
