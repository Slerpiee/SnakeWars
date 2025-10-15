package models

import "math"

type Point struct{
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Speed struct{
	Dx int `json:"dx"`
	Dy int	`json:"dy"`
}

func distance(p1, p2 Point) float64{
	return math.Sqrt(math.Pow((p1.X-p2.X), 2) + math.Pow((p1.Y - p2.Y), 2))
}

func direction(from, to Point) float64{
	return math.Atan2((to.Y-from.Y), (to.X - from.X))
}


type Segment struct{
	Start Point
	End Point
	Length float64
	Direction float64 //radians
}

const MIN_SEGMENT_LEN = 5.0 //Минимальная длина змейки 

type User struct{
	ID string `json:"id"`
	Username string `json:"username"`
	Room_id string
}

type PlayerCosmetics struct{}

type PlayerState struct{
	isAlive bool
	isReady bool
}

type Snake struct{
	ID string
	Head Point
	Speed Speed
	Skin *PlayerCosmetics
	State PlayerState
	Body []Segment //очередь
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

func (s *Snake) Move(p Point){
	var eps_rad = 0.1 //Сравнивать float64 некорректно, будем считать, что x и y равны, если |x-y| <= eps (ПОДОБРАТЬ ОПТИМАЛЬНОЕ ЗНАЧЕНИЕ ВО ВРЕМЯ ДЕБАГА)
	if (len(s.Body) == 0){
		s.Body = append(s.Body, Segment{s.Head, p, distance(s.Head, p), direction(s.Head, p)})
		s.SegmentCounter = 1
	} else{
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

//func NewSnake(id string, startPoint Point, color string) *Snake{
	// StartBody := make([]Point{startPoint})
	// return &Snake{
	// 	ID: id,
	// 	Head: startPoint,
	// 	Speed: Speed{
	// 		Dx: 0,
	// 		Dy: 0
	// 	},
	// 	Body: StartBody,
	// 	Length: 1, //Бессмысленное поле
	// 	Color: color,
	// 	IsAlive: true,
	// }
//}




