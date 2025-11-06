package game

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type UserState int

const (
	UserStateNotReady UserState = 0
	UserStateReady    UserState = 1
	UserStateKilled   UserState = 2
	UserStateExit     UserState = 3
)

type User struct {
	ID       string
	Conn     *websocket.Conn
	Name     string
	RoomID   string
	Snake    *Snake
	State    UserState
	SendChan chan any //Канал для отправки сообщений клиенту
	mutex    sync.RWMutex
	LastPing time.Time
}

//Слушатель SendChan для отправки сообщений пользователю через websocket

func (u *User) Close() {
	ch := u.SendChan
	if ch != nil {
		close(ch)
	}
}
