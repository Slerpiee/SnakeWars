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
	Body  []Segment //очередь
	SegmentCounter int
}

func (s *Snake) Add_Segment(seg Segment) int{
	s.Body = append(s.Body, seg)
	s.SegmentCounter += 1
	return s.SegmentCounter
}

func (s *Snake) Remove_Segment(seg Segment) int{
	s.SegmentCounter > 0
}

//База:
func NewSnake(id string, startPoint Point, color string) *Snake {
	StartBody := []Segment{
		{
			Start:     startPoint,
			End:       startPoint,
			Length:    0,
			Direction: 0,
		},
	}
	return &Snake{
		ID:    id,
		Head:  startPoint,
		Speed: Speed{Dx: 1, Dy: 0},
		Skin:  PlayerCosmetics{Color: color},
		State: PlayerState{isAlive: true, isReady: false},
		Stats: PlayerStats{
			Wins:        0,
			Losses:      0,
			GamesPlayed: 0,
			MaxLength:   0,
		},
		Body: StartBody,
	}
}

func (s *Snake) Grow(amount float64){
	//добавление в кольцевой буффер
}

func (s *Snake) GetHead() Point{
	point := Point{s.Head.X, s.Head.Y}
	return point
}


//Движение:
func (s *Snake) Move(p Point) {

	if len(s.Body) == 0 {
		s.Body = append(s.Body, Segment{s.Head, p, distance(s.Head, p), direction(s.Head, p)})
		s.SegmentCounter = 1
	} else {
		head_segment := &s.Body[len(s.Body)-1]
		dist := distance(head_segment.End, p)
		new_direction := direction(head_segment.End, p)
		if math.Abs(head_segment.Direction - new_direction) < eps_rad {
			head_segment.End = p
			head_segment.Length += dist
		} else{
			new_seg := Segment{s.Head, p, dist, new_direction}
			
		}

	}
}
