package storage

import "time"

type Notification struct {
	// ID события;
	ID string
	// Заголовок события;
	Title string
	// Дата события;
	DateTime time.Time
	// Пользователь, которому отправлять.
	UserID int64
}
