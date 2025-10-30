package game

import (
	"github.com/gorilla/websocket"
	"time"
	"sync"
)

type UserState int

const (
    UserStateNotReady UserState = 0
    UserStateReady UserState = 1
    UserStateKilled UserState = 2
	UserStateExit UserState = 3
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
	ch := u.SendChan
	if ch != nil{
		close(ch)
	}
}

