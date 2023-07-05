package main

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
	len  int
	head *ListItem
	end  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.end
}

func (l *list) PushFront(v interface{}) *ListItem {
	prevHead := l.head
	l.head = &ListItem{v, prevHead, nil}

	if l.len == 0 {
		l.end = l.head
	}

	if prevHead != nil {
		prevHead.Prev = l.head
	}
	l.len++

	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	prevEnd := l.end
	if l.len == 0 {
		l.head = &ListItem{v, prevEnd, nil}
		l.end = l.head
		l.len++
	} else {
		l.end = &ListItem{v, nil, prevEnd}
		prevEnd.Next = l.end
		l.len++
	}
	if prevEnd == l.head {
		l.head.Next = l.end
	}
	return l.end
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	switch {
	case i == l.end:
		l.end = i.Prev
		l.end.Next = nil
	case i == l.head:
		l.head = nil
		l.head.Prev = nil
	case i != l.head:
		i.Prev.Next = i.Next
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.head {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}
