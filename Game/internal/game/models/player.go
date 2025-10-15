package models

type Point struct{
	x int `json:"x"`
	y int `json:"y"`
}

type Speed struct{
	Dx int `json:"dx"`
	Dy int	`json:"dy"`
}

type User struct{
	ID string `json:"id"`
	Username string `json:"username"`
}

type Snake struct{
	ID string
	Room_id string
	Head Point
	Speed Speed
	Length int
	Color string
	IsAlive bool
	Body []Point //очередь
}

func NewSnake(id string, startPoint Point, color string) *Snake{
	StartBody := make([]Point{startPoint})
	return &Snake{
		ID: id,
		Head: startPoint,
		Speed: Speed{
			Dx: 1,
			Dy: 0
		},
		Body: StartBody,
		Length: 1,
		Color: color,
		IsAlive: true,
	}
}




