package game

import (
	"sync"
	"time"

	"github.com/gin/websocket"
)

type RoomState struct{}

type RoomStats struct{}


type Room struct{
	ID string //Room ID

	
	Clients sync.Map//map[string]*Client //ID: *Client


	Room_mutex sync.RWMutex
	State RoomState
	Stats RoomStats

	Ticker *time.Ticker //update per second 

}
