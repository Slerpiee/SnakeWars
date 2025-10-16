package models

import (
	"math"
)



type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}


type Speed struct {
	Dx float64 `json:"dx"`
	Dy float64 `json:"dy"`
}



func distance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow((p1.X-p2.X), 2) + math.Pow((p1.Y-p2.Y), 2))
}

func direction(from, to Point) float64 {
	return math.Atan2((to.Y - from.Y), (to.X - from.X))
}

type Segment struct {
	Start     Point
	End       Point
	Length    float64
	Direction float64 //radians

}

const MIN_SEGMENT_LEN = 5.0 //Минимальная длина змейки 
const MAX_SEGMENTS = 100

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Room_id  string
}

type PlayerCosmetics struct {
	Color string
}

type PlayerState struct {
	isAlive bool
	isReady bool
}

type PlayerStats struct {
	Wins        int
	Losses      int
	GamesPlayed int
	MaxLength   int
}

type Snake struct {
	ID    string
	Head  Point
	Speed Speed
	Skin  PlayerCosmetics
	State PlayerState
	Stats PlayerStats
	Body  *RingBuffer[Segment] 
}


// func (s *Snake) Add_Segment(seg Segment) int{
// 	s.Body = append(s.Body, seg)
// 	s.SegmentCounter += 1
// 	return s.SegmentCounter
// }

// func (s *Snake) Remove_Segment(seg Segment) int{
// 	s.SegmentCounter > 0
// }

//База:
func NewSnake(id string, startPoint Point, color string) *Snake {
	StartBody := NewRingBuffer[Segment](MAX_SEGMENTS)
	return &Snake{
		ID:    id,
		Head:  startPoint,
		Speed: Speed{},
		Skin:  PlayerCosmetics{Color: color},
		State: PlayerState{},
		Stats: PlayerStats{
			Wins:        0,
			Losses:      0,
			GamesPlayed: 0,
			MaxLength:   0,
		},
		Body: StartBody,
	}
}


func (s *Snake) GetHead() *Point{
	point := &Point{s.Head.X, s.Head.Y}
	return point
}


//Движение: p1 = {X, Y} -> p2 = {X+dx*T, Y + dy*t}
func (s *Snake) Move(p Point, grow bool) { 
	eps_rad := 0.01 //Придется подбирать методом подбора
	if s.Body.size == 0 {
		s.Body.Push(Segment{s.Head, p, distance(s.Head, p), direction(s.Head, p)})
	} else {
		head_segment := s.Body.Peek().(Segment)
		dist := distance(head_segment.End, p)
		new_direction := direction(head_segment.End, p)
		if math.Abs(head_segment.Direction - new_direction) <= eps_rad { //Незначительное отклонение от направления головного куска, просто продлеваем головной кусок
			head_segment.End = p
			head_segment.Length += dist
		} else{ //Значительное отклонение, создаем новый кусок
			if s.Body.IsFull(){ 
				s.Body.Pop()
			}
			new_seg := Segment{s.Head, p, dist, new_direction}
			s.Body.Push(new_seg)
		}
		if !grow{
			s.shrink(dist)
		}
	}

}
 
func (s *Snake) shrink(dist_to_remove float64){
	if s.Body.size == 0{
		return
	}
	for dist_to_remove > 0 && s.Body.size > 0 {
		tail_seg := &s.Body.buffer[s.Body.head]
		if tail_seg.Length <= dist_to_remove{ 
			dist_to_remove -= tail_seg.Length
			s.Body.Pop()
		} else{
			ratio := dist_to_remove / tail_seg.Length
			prev_length := tail_seg.Length
			dx := tail_seg.End.X - tail_seg.Start.X
			dy := tail_seg.End.Y - tail_seg.Start.Y
			tail_seg.Start.X = tail_seg.Start.X + dx * ratio
			tail_seg.Start.Y = tail_seg.Start.Y + dy * ratio
			tail_seg.Length = prev_length - dist_to_remove
			return 
		}

	}

}

