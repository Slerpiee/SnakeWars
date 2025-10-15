package models

type Node struct{
	Val Segment
	next *Node
}

type Queue struct{
	last *Node
	first *Node
	len int
}

func NewQueue() *Queue {
    return &Queue{}
}

func (q *Queue) Push(seg Segment){
	new_node := &Node{Val: seg, next: nil}
	if (q.last == nil){
		q.first = new_node
		q.last = new_node
	}else{
		q.last.next = new_node
		q.last = new_node
	}
	q.len ++
}

func (q *Queue) Pop() Segment{
	if q.first == nil{
		panic("Queue is empty")
	}
	val := q.first.Val
	q.first = q.first.next
	if q.first == nil{
		q.last = nil
	}
	q.len --
	return val
}

func (q *Queue) IsEmpty() bool{
	return q.first == nil
}

func (q *Queue) Len() int {
    return q.len
}

func (q *Queue) Peek() Segment {
    if q.first == nil {
        panic("Queue is empty")
    }
    return q.first.Val
}