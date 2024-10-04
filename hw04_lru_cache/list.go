package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	result := &ListItem{
		Value: v,
	}
	l.len++
	front := l.front
	if front != nil {
		result.Next = front
		front.Prev = result
	} else if l.back == nil {
		l.back = result
	}
	l.front = result
	return result
}

func (l *list) PushBack(v interface{}) *ListItem {
	result := &ListItem{
		Value: v,
	}
	l.len++
	back := l.back
	if back != nil {
		result.Prev = back
		back.Next = result
	} else if l.front == nil {
		l.front = result
	}
	l.back = result
	return result
}

func (l *list) Remove(i *ListItem) {
	next := i.Next
	prev := i.Prev
	if next != nil {
		next.Prev = prev
	} else if i == l.back {
		l.back = prev
	}
	if prev != nil {
		prev.Next = next
	} else if i == l.front {
		l.front = next
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}
	next := i.Next
	prev := i.Prev
	if next != nil {
		next.Prev = prev
	}
	if prev != nil {
		prev.Next = next
	}
	if i == l.back {
		l.back = prev
	}
	i.Next = l.front
	i.Prev = nil
	l.front = i
}

func NewList() List {
	return new(list)
}
