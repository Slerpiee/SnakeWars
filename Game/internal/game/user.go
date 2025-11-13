package game

import (
	"sync"
	"time"
	"log"
	"strings"
	"errors"
	"io"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
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
	RoomID   *Room
	Snake    *Snake
	State    UserState
	Chan chan Message //Канал для отправки сообщений клиенту
	mutex    sync.RWMutex
	LastPing time.Time
}



//Слушатель SendChan для отправки сообщений пользователю через websocket

func CreateUser(conn *websocket.Conn) *User{
	return &User{
		ID:uuid.NewString(),
		Conn: conn, 
	}
}


func (u *User) ReadPump(){
	defer func() {
        u.Conn.Close()
        u.CloseChan()
        log.Printf("User %s read pump stopped", u.Name)
    }()
	wsMessages := make(chan Message, 100) //Канал с сообщениями для клиента
	go func() {
		defer close(wsMessages)
        for {
            var message Message
            err := u.Conn.ReadJSON(&message)
            if err != nil {
                if isRecoverableError(err){
					continue
				}
                close(wsMessages)
                return
            }
			wsMessages <- message
        }
    }()
	for {
		select{
		case m, ok := <- u.Chan: //messages like player dead, moved etc
			if !ok{
				return 
			}
			log.Printf("User %s get message from server channel: %s", u.Name, m.ID)
		case m, ok := <- wsMessages: //messages like user pressed forward, backward etc
			if !ok{
				log.Printf("Closed ws channel for user: %s", u.Name)
				return
			}
			log.Printf("User %s send message %s via websocket to server: ", u.Name, m.ID )
		}
	}
}

func (u *User) UpdatePos(){
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.Snake.Move(time.Since(u.LastPing), false)
}


func (u *User) CloseChan() {
	if u.Chan != nil {
		close(u.Chan)
	}
}



func isRecoverableError(err error) bool {

    if strings.Contains(err.Error(), "json:") ||
       strings.Contains(err.Error(), "unmarshal") ||
       strings.Contains(err.Error(), "invalid character") {
        return true
    }
    

    if websocket.IsCloseError(err, 
        websocket.CloseNormalClosure,
        websocket.CloseGoingAway,
        websocket.CloseAbnormalClosure) {
        return false
    }
    

    if errors.Is(err, io.ErrUnexpectedEOF) {
        return true
    }
    
    return false
}

