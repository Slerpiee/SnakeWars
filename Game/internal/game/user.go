package game

import (
	"github.com/gorilla/websocket"
	"time"
	"sync"
)

type UserState int

const (
    USERSTATE_NOTREADY UserState = 0
    USERSTATE_READY UserState = 1
    USERSTATE_KILLED UserState = 2
	USERSTATE_EXIT UserState = 3
)


type User struct {
	ID    string
	Conn  *websocket.Conn
	Name string
	RoomID string
	Snake *Snake
	State UserState
	SendChan  chan any //Канал для отправки сообщений клиенту
	mutex    sync.RWMutex
    LastPing time.Time
}

func (u *User) Close(){
	if u.SendChan != nil{
		close(u.SendChan)
		u.SendChan = nil
	}
}

