package storage

import "errors"

var (
	ErrEventNotFound = errors.New("событие не найдено")
	ErrDateBusy      = errors.New("данное время уже занято другим событием")
)
