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
	listItems map[*ListItem]struct{}
	front     *ListItem
	back      *ListItem
}

func (l *list) Len() int {
	return len(l.listItems)
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
	l.listItems[result] = struct{}{}
	front := l.front
	if front != nil {
		result.Next = front
		front.Prev = result
	}
	l.front = result
	if l.back == nil {
		l.back = result
	}
	return result
}

func (l *list) PushBack(v interface{}) *ListItem {
	result := &ListItem{
		Value: v,
	}
	l.listItems[result] = struct{}{}
	back := l.back
	if back != nil {
		result.Prev = back
		back.Next = result
	}
	l.back = result
	if l.front == nil {
		l.front = result
	}
	return result
}

func (l *list) Remove(i *ListItem) {
	next := i.Next
	prev := i.Prev
	if next != nil {
		next.Prev = prev
	}
	if prev != nil {
		prev.Next = next
	}
	if i == l.front {
		l.front = next
	}
	if i == l.back {
		l.back = prev
	}
	delete(l.listItems, i)
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
	result := new(list)
	result.listItems = make(map[*ListItem]struct{})
	return result
}
