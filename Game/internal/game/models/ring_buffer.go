package models


type RingBuffer[T any] struct {
    buffer   []T
    head     int 
    tail     int 
    size     int 
    capacity int 
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
    return &RingBuffer[T]{
        buffer:   make([]T, capacity),
        capacity: capacity,
        head:     0,
        tail:     0,
        size:     0,
    }
}

func (r *RingBuffer[T]) Push(item T) bool {
    if r.size == r.capacity {
        return false 
    }
    
    r.buffer[r.tail] = item
    r.tail = (r.tail + 1) % r.capacity  //th -> h-t		
    r.size++
    return true
}

func (r *RingBuffer[T]) Pop() T {
    var zero T
    if r.size == 0 {
        return zero
    }
    
    item := r.buffer[r.head]
    r.head = (r.head + 1) % r.capacity
    r.size--
    return item
}


func (r *RingBuffer[T]) IsEmpty() bool {
    return r.size == 0
}


func (r *RingBuffer[T]) IsFull() bool {
    return r.size == r.capacity
}


func (r *RingBuffer[T]) Size() int {
    return r.size
}


func (r *RingBuffer[T]) Peek() interface{} {
    if r.IsEmpty() {
        return nil
    }
    return r.buffer[r.tail]
}
