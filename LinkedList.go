package main

type Node struct {
	info Pixel
	next *Node
}

type LinkedList struct {
	head *Node
}

func (l *LinkedList) Append(p Pixel) {
	list := &Node{info: p, next: nil}
	if l.head == nil {
		l.head = list
	} else {
		p := l.head
		for p.next != nil {
			p = p.next
		}
		p.next = list
	}
}
