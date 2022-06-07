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
	firstItem *ListItem
	listItem  *ListItem
	len       int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.listItem
}

// PushFront move item to head list if list have this item.
func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.firstItem == nil {
		l.firstItem = newItem
		l.listItem = newItem
	} else {
		newItem.Next = l.firstItem
		l.firstItem.Prev = newItem
		l.firstItem = newItem
	}
	l.len++
	return newItem
}

// PushBack move item to tail list if list have this item.
func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.listItem == nil {
		return l.PushFront(v)
	}
	newItem.Prev = l.listItem
	l.listItem.Next = newItem
	l.listItem = newItem
	l.len++
	return newItem
}

// Remove item if list have this item.
func (l *list) Remove(i *ListItem) {
	if l.len == 0 {
		return
	}
	if i.Prev == nil {
		l.firstItem = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.listItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.len--
}

// MoveToFront move item to head list if list have this item.
func (l *list) MoveToFront(i *ListItem) {
	if l.firstItem == i {
		return
	}
	if i.Next != nil && i.Prev != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	if i.Next == nil {
		l.listItem = i.Prev
		i.Prev.Next = nil
	}

	currentHead := l.firstItem
	l.firstItem = i
	i.Prev = nil
	i.Next = currentHead
	currentHead.Prev = l.firstItem
}

func NewList() List {
	return new(list)
}
