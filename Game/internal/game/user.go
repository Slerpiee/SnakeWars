package game

import (
	"github.com/gorilla/websocket"
	"time"
	"sync"
)

type User struct {
	ID    string
	Conn  *websocket.Conn
	Snake *Snake
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

