package dto

type Event struct {
	// ID - уникальный идентификатор события (можно воспользоваться UUID);
	ID string `json:"id"`
	// Заголовок - короткий текст;
	Title string `json:"title"`
	// Дата и время события;
	DateTime string `json:"dateTime"`
	// Длительность события (или дата и время окончания);
	Duration string `json:"duration"`
	// Описание события - длинный текст, опционально;
	// ID пользователя, владельца события;
	UserID int64 `json:"userId"`
	// За сколько времени высылать уведомление, опционально.
	NotificationDuration string `json:"notificationDuration"`
}
