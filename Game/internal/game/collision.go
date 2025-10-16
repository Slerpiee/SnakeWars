package game
const SNAKE_WIDTH = 5

func CheckHeadsCollision(s1, s2 *Snake) bool{
	if s1.ID == s2.ID{
		return CheckSelfCollision(s1)
	}
	head1 := s1.Head
	head2 := s2.Head
	dist_heads := distance(head1, head2)
	return dist_heads <= SNAKE_WIDTH*2
}

func CheckSelfCollision(s1 *Snake) bool{
	head := s1.Head
	size := s1.Body.Size()
	body := s1.Body
	if size <= 3{
		return false
	}

	for i := size - 1; i >= 3; i --{ //от конца до 3 сегмента 
		cur_seg := *(body.Get(i).(*Segment))
		dist := DistancePointSegment(head, cur_seg)
		if dist <= SNAKE_WIDTH*2{
			return true
		}
	}
	return false
}

func CheckBodyCollision(s1, s2 *Snake) bool{ // проверяем, что s1 сталкивается ебалом с телом s2. Чтобы проверить обратное, просто поменяем порядок аргументов
	if s2.Body.Size() <= 1{
		return false
	}
	head := s1.Head
	size := s2.Body.Size()
	body := s2.Body

	for i := size - 1; i > 0; i --{ //от конца до 2 сегмента 
		cur_seg := *(body.Get(i).(*Segment)) 
		dist := DistancePointSegment(head, cur_seg)
		if dist <= SNAKE_WIDTH*2{
			return true
		}
	}
	return false
}

func CheckAllCollisions(s1, s2 *Snake) string {
	//Проверяем самопересечения
    if CheckSelfCollision(s1) {
        return s2.ID
    }
    if CheckSelfCollision(s2) {
        return s1.ID
    }
    
    //Проверяем столкновение голов
    if CheckHeadsCollision(s1, s2) {
        return "tie"
    }
    
    //Проверяем столкновения с телами
    s1HitsS2Body := CheckBodyCollision(s1, s2)
    s2HitsS1Body := CheckBodyCollision(s2, s1)
    
    if s1HitsS2Body && s2HitsS1Body {
        return "tie"
    }
    if s1HitsS2Body {
        return s2.ID
    }
    if s2HitsS1Body {
        return s1.ID
    }
    
    return "no collisions"
}