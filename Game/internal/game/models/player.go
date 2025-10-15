package models

type Point struct{
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Speed struct{
	Dx int `json:"dx"`
	Dy int	`json:"dy"`
}




type Segment struct{
	Start Point
	End Point
	Length float64
	Direction Speed 

}

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
}

func (s *Snake) Move(p Point){
	
	if (len(s.Body) == 0){
		s.Body = append(s.Body, Segment{p, p, 0, Speed{0, 1}})
	} else{
		
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




