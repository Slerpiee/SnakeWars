package game

import "math"

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


func DistancePointSegment(p Point, seg Segment) float64 {
	// Вектор отрезка AB
	abX := seg.End.X - seg.Start.X
	abY := seg.End.Y - seg.Start.Y
	
	// Вектор от точки A до точки P
	apX := p.X - seg.Start.X
	apY := p.Y - seg.Start.Y
	
	// Скалярное произведение AP · AB
	dotAPAB := apX*abX + apY*abY
	
	// Квадрат длины отрезка AB
	abLengthSquared := abX*abX + abY*abY
	
	// Если отрезок вырожден (начало и конец совпадают), 
	// возвращаем расстояние до любой из точек
	if abLengthSquared == 0 {
		return distance(p, seg.Start)
	}
	
	// Вычисляем параметр t - положение проекции на прямой
	t := dotAPAB / abLengthSquared
	
	// Ограничиваем t диапазоном [0, 1] для работы с отрезком, а не с прямой
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	
	// Находим координаты ближайшей точки на отрезке
	closestX := seg.Start.X + t*abX
	closestY := seg.Start.Y + t*abY
	
	// Возвращаем расстояние до ближайшей точки
	return distance(p, Point{X: closestX, Y: closestY})
}

func DistanceSegmentToSegment(seg1, seg2 Segment) float64 {
    // Векторы отрезков
    uX := seg1.End.X - seg1.Start.X
    uY := seg1.End.Y - seg1.Start.Y
    vX := seg2.End.X - seg2.Start.X
    vY := seg2.End.Y - seg2.Start.Y
    wX := seg1.Start.X - seg2.Start.X
    wY := seg1.Start.Y - seg2.Start.Y
    
    a := uX*uX + uY*uY // всегда >= 0
    b := uX*vX + uY*vY
    c := vX*vX + vY*vY // всегда >= 0
    d := uX*wX + uY*wY
    e := vX*wX + vY*wY
    
    denom := a*c - b*b // всегда >= 0
    
    var sN, sD, tN, tD float64
    
    // Вычисляем параметры для ближайших точек на бесконечных прямых
    if denom < 1e-10 {
        // Отрезки параллельны - обрабатываем как вырожденный случай
        sN = 0.0
        sD = 1.0
        tN = e
        tD = c
    } else {
        sN = b*e - c*d
        tN = a*e - b*d
        if sN < 0.0 {
            sN = 0.0
            tN = e
            tD = c
        } else if sN > denom {
            sN = denom
            tN = e + b
            tD = c
        } else {
            sD = denom
            tD = denom
        }
    }
    
    if tN < 0.0 {
        tN = 0.0
        // Пересчитываем sN для t = 0
        if -d < 0.0 {
            sN = 0.0
        } else if -d > a {
            sN = sD
        } else {
            sN = -d
            sD = a
        }
    } else if tN > tD {
        tN = tD
        // Пересчитываем sN для t = 1
        if (-d + b) < 0.0 {
            sN = 0.0
        } else if (-d + b) > a {
            sN = sD
        } else {
            sN = -d + b
            sD = a
        }
    } else {
        sD = denom
        tD = denom
    }
    
    // Вычисляем параметры s и t
    var s, t float64
    if math.Abs(sN) < 1e-10 {
        s = 0.0
    } else {
        s = sN / sD
    }
    
    if math.Abs(tN) < 1e-10 {
        t = 0.0
    } else {
        t = tN / tD
    }
    
    // Вычисляем ближайшие точки
    closest1 := Point{
        X: seg1.Start.X + s*uX,
        Y: seg1.Start.Y + s*uY,
    }
    
    closest2 := Point{
        X: seg2.Start.X + t*vX,
        Y: seg2.Start.Y + t*vY,
    }
    
    // Возвращаем расстояние между ближайшими точками
    return distance(closest1, closest2)
}



func SegmentsCollide(seg1, seg2 Segment, tolerance float64) bool {
    dist := DistanceSegmentToSegment(seg1, seg2)
    return dist <= tolerance
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

func CheckCollision(s1, s2 *Snake) (string)  { // func (s1, s22) -> (pointer_to_winner, isColided)
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