package storage

import (
	"errors"
	"fmt"
)

var (
	ErrBusiness      = errors.New("business error")
	ErrEventNotFound = fmt.Errorf("%w : событие не найдено", ErrBusiness)
	ErrDateBusy      = fmt.Errorf("%w : данное время уже занято другим событием", ErrBusiness)
)
