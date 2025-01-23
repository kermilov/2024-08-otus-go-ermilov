package storage

import "time"

type Event struct {
	// ID - уникальный идентификатор события (можно воспользоваться UUID);
	ID string
	// Заголовок - короткий текст;
	Title string
	// Дата и время события;
	DateTime time.Time
	// Длительность события (или дата и время окончания);
	Duration time.Duration
	// Описание события - длинный текст, опционально;
	// ID пользователя, владельца события;
	UserID int64
	// За сколько времени высылать уведомление, опционально.
	NotificationDuration time.Duration
}
