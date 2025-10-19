package game

import "time"

import "math"

const SNAKE_WIDTH = 3


// DistanceSegmentToSegment возвращает минимальное расстояние между двумя отрезками
func DistanceSegmentToSegment(seg1, seg2 Segment) float64 {
    // Векторы отрезков
    uX := seg1.End.X - seg1.Start.X
    uY := seg1.End.Y - seg1.Start.Y
    vX := seg2.End.X - seg2.Start.X
    vY := seg2.End.Y - seg2.Start.Y
    wX := seg1.Start.X - seg2.Start.X
    wY := seg1.Start.Y - seg2.Start.Y

    a := uX*uX + uY*uY // длина seg1 в квадрате, всегда >= 0
    b := uX*vX + uY*vY
    c := vX*vX + vY*vY // длина seg2 в квадрате, всегда >= 0
    d := uX*wX + uY*wY
    e := vX*wX + vY*wY

    denom := a*c - b*b // determinant, всегда >= 0

    var sN, sD, tN, tD float64
    sD = denom
    tD = denom

    // Вычисляем параметры для ближайших точек на бесконечных прямых
    if denom < 1e-10 {
        // Отрезки параллельны
        sN = 0.0
        sD = 1.0
        tN = e
        tD = c
    } else {
        sN = b*e - c*d
        tN = a*e - b*d
        
        if sN < 0.0 {
            // sN < 0 => the s=0 edge is visible
            sN = 0.0
            tN = e
            tD = c
        } else if sN > sD {
            // sN > sD => the s=1 edge is visible
            sN = sD
            tN = e + b
            tD = c
        }
    }

    if tN < 0.0 {
        // tN < 0 => the t=0 edge is visible
        tN = 0.0
        // recompute sN for this edge
        if -d < 0.0 {
            sN = 0.0
        } else if -d > a {
            sN = sD
        } else {
            sN = -d
            sD = a
        }
    } else if tN > tD {
        // tN > tD => the t=1 edge is visible
        tN = tD
        // recompute sN for this edge
        if (-d + b) < 0.0 {
            sN = 0.0
        } else if (-d + b) > a {
            sN = sD
        } else {
            sN = -d + b
            sD = a
        }
    }

    // finally do the division to get sc and tc
    sc := 0.0
    if math.Abs(sN) < 1e-10 {
        sc = 0.0
    } else {
        sc = sN / sD
    }

    tc := 0.0
    if math.Abs(tN) < 1e-10 {
        tc = 0.0
    } else {
        tc = tN / tD
    }

    // get the difference of the two closest points
    dX := wX + (sc * uX) - (tc * vX)
    dY := wY + (sc * uY) - (tc * vY)

    return math.Sqrt(dX*dX + dY*dY)
}

func SegmentsCollide(seg1, seg2 Segment, tolerance float64) bool {
    // Вычисляем расстояние между сегментами
    dist := DistanceSegmentToSegment(seg1, seg2)
    
    // Если расстояние меньше или равно tolerance - сегменты сталкиваются
    return dist <= tolerance
}

// ContinuousCollision проверяет столкновение между движущейся головой и сегментом тела
// за промежуток времени deltaTime
func ContinuousCollision(headStart, headEnd Point, bodySeg Segment, snakeWidth float64) bool {
    // Создаем сегмент движения головы
    headMovement := Segment{
        Start: headStart,
        End:   headEnd,
        Length: distance(headStart, headEnd),
    }
    
    // Проверяем столкновение сегмента движения головы с сегментом тела
    return SegmentsCollide(headMovement, bodySeg, snakeWidth*2)
}

// CheckContinuousBodyCollision проверяет столкновение головы s1 с телом s2 с учетом движения
func CheckContinuousBodyCollision(s1, s2 *Snake, deltaTime time.Duration) bool {
    if s2.Body.Size() <= 1 {
        return false
    }
    
    // Вычисляем предыдущую позицию головы s1
    dt := deltaTime.Seconds()
    prevHead := Point{
        X: s1.Head.X - s1.Speed.Dx*dt,
        Y: s1.Head.Y - s1.Speed.Dy*dt,
    }
    
    // Проверяем столкновение с каждым сегментом тела s2
    for i := s2.Body.Size() - 1; i > 0; i-- {
        bodySeg := *(s2.Body.Get(i).(*Segment))
        if ContinuousCollision(prevHead, s1.Head, bodySeg, SNAKE_WIDTH) {
            return true
        }
    }
    
    return false
}

// CheckContinuousHeadsCollision проверяет столкновение голов с учетом движения
func CheckContinuousHeadsCollision(s1, s2 *Snake, deltaTime time.Duration) bool {
    if s1.ID == s2.ID {
        return CheckContinuousSelfCollision(s1, deltaTime)
    }
    
    dt := deltaTime.Seconds()
    prevHead1 := Point{
        X: s1.Head.X - s1.Speed.Dx*dt,
        Y: s1.Head.Y - s1.Speed.Dy*dt,
    }
    prevHead2 := Point{
        X: s2.Head.X - s2.Speed.Dx*dt,
        Y: s2.Head.Y - s2.Speed.Dy*dt,
    }
    
    // Проверяем столкновение сегментов движения голов
    headMovement1 := Segment{Start: prevHead1, End: s1.Head}
    headMovement2 := Segment{Start: prevHead2, End: s2.Head}
    
    return SegmentsCollide(headMovement1, headMovement2, SNAKE_WIDTH*2)
}

// CheckContinuousSelfCollision проверяет самопересечение с учетом движения
func CheckContinuousSelfCollision(s1 *Snake, deltaTime time.Duration) bool {
    if s1.Body.Size() <= 3 {
        return false
    }
    
    dt := deltaTime.Seconds()
    prevHead := Point{
        X: s1.Head.X - s1.Speed.Dx*dt,
        Y: s1.Head.Y - s1.Speed.Dy*dt,
    }
    
    // Проверяем столкновение с сегментами тела (кроме последних 2, чтобы избежать
    // ложных срабатываний на соседних сегментах)
    for i := s1.Body.Size() - 1; i >= 3; i-- {
        bodySeg := *(s1.Body.Get(i).(*Segment))
        if ContinuousCollision(prevHead, s1.Head, bodySeg, SNAKE_WIDTH) {
            return true
        }
    }
    
    return false
}

// CheckContinuousCollision основная функция проверки коллизий с непрерывным обнаружением
func CheckContinuousCollision(s1, s2 *Snake, deltaTime time.Duration) (*Snake, bool) { //snakeId, isTie
    // Проверяем самопересечения
    if CheckContinuousSelfCollision(s1, deltaTime) {
        return s2, false
    }
    if CheckContinuousSelfCollision(s2, deltaTime) {
        return s1, false
    }
    
    // Проверяем столкновение голов
    if CheckContinuousHeadsCollision(s1, s2, deltaTime) {
        return s1, true
    }
    
    // Проверяем столкновения с телами
    s1HitsS2Body := CheckContinuousBodyCollision(s1, s2, deltaTime)
    s2HitsS1Body := CheckContinuousBodyCollision(s2, s1, deltaTime)
    
    if s1HitsS2Body && s2HitsS1Body {
        return s1, true
    }
    if s1HitsS2Body {
        return s2, false
    }
    if s2HitsS1Body {
        return s1, false
    }
    
    return nil, false
}