package utils


type RingBuffer[T any] struct {
    Buffer   []T
    head     int 
    tail     int 
    size     int 
    Capacity int 
}

func NewRingBuffer[T any](Capacity int) *RingBuffer[T] {
    return &RingBuffer[T]{
        Buffer:   make([]T, Capacity),
        Capacity: Capacity,
        head:     0,
        tail:     0,
        size:     0,
    }
}

func (r *RingBuffer[T]) Push(item T) bool {
    if r.size == r.Capacity {
        return false 
    }
    
    r.Buffer[r.tail] = item
    r.tail = (r.tail + 1) % r.Capacity  //th -> h-t		
    r.size++
    return true
}

func (r *RingBuffer[T]) Pop() T {
    var zero T
    if r.size == 0 {
        return zero
    }
    
    item := r.Buffer[r.head]
	r.Buffer[r.head] = zero //для сборщика мусора
    r.head = (r.head + 1) % r.Capacity
    r.size--
    return item
}


func (r *RingBuffer[T]) IsEmpty() bool {
    return r.size == 0
}


func (r *RingBuffer[T]) IsFull() bool {
    return r.size == r.Capacity
}


func (r *RingBuffer[T]) Size() int {
    return r.size
}


func (r *RingBuffer[T]) Peek() any { //Возвращает указатель на последний добавленный элеменет
    if r.IsEmpty() {
        return nil
    }
    return &r.Buffer[r.tail]
}

func (r *RingBuffer[T]) Tail() any { //Возвращает указатель на элемент, который выйдет из очереди первым
    if r.IsEmpty(){
        return nil
    }
    return &r.Buffer[r.head]
}


