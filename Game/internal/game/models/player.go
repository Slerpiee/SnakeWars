package models

import "math"

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Speed struct {
	Dx int `json:"dx"`
	Dy int `json:"dy"`
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

type User struct{
	ID string `json:"id"`
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
	Body []Segment //очередь
}

func (s *Snake) Move(p Point){
	
	if (len(s.Body) == 0){
		s.Body = append(s.Body, Segment{s.Head, p, distance(s.Head, p), direction(s.Head, p)})
	} else{
		//dist := s.Body[len(s.Body)-1]
		
	}
}
