package game

import (
	"sync"
	"time"
	"log"
	"github.com/gorilla/websocket"
)

type RoomState struct{}

type RoomStats struct{}

type RoomMessage struct {
    Message  any //Дописать структуру для сообщений
}


type Player struct {
	ID    string
	Conn  *websocket.Conn
	Snake *Snake
	SendChan  chan any //Канал для отправки сообщений клиенту
	mutex    sync.RWMutex
    LastPing time.Time
}

type Room struct{
	ID string //Room ID

	
	players sync.Map  //ID: *Player


	room_mutex sync.RWMutex
	State RoomState
	Stats RoomStats

	roomInput chan RoomMessage
	//RoomChan chan
	ticker *time.Ticker //update per second 

}


func (room *Room) addPlayer(player *Player) {
    room.players.Store(player.ID, player)
}


func (room *Room) getPlayer(id string) *Player {
	value, ok := room.players.Load(id)
    if ok {
        return value.(*Player)
    }
    return nil
}


func (room *Room) removePlayer(id string){
    room.players.Delete(id)
}


func (room *Room) UpdatePlayer(id string) {
    if player := room.getPlayer(id); player != nil {
        player.mutex.Lock() 
        defer player.mutex.Unlock()
		player.Snake.Move(time.Now().Sub(player.LastPing), false)
    }
}

func (room *Room) Broadcast(message interface{}) {
    room.players.Range(func(key, value interface{}) bool {
        player := value.(*Player)
        
        select {
        	case player.SendChan <- message:
        default:
            log.Printf("Player %s is lagging", player.ID) //Пользователь не забирает сообщения из канала
        }
        return true
    })
}

func (room *Room) startInputProcessor() {
    room.roomInput = make(chan RoomMessage, 100)
    
    go func() {
        for range room.roomInput { //Апдейты типо пользователь нажал кнопку управления, кто-то ебнулся об стенку и сдох и тд
            //handle room update
			log.Printf("Message for room %s: %", room.ID)
        }
    }()
}

func (room *Room) startGameLoop() {
    room.ticker = time.NewTicker(16 * time.Millisecond) // 60 FPS
    
    go func() {
        for range room.ticker.C {
            room.gameTick()
        }
    }()
}


func (room *Room) gameTick() {
    room.room_mutex.Lock()
	//r.updateSnakes()
	//if checkCollisions -> kill
	//... BroadCast informatoin about kill or sum
    defer room.room_mutex.Unlock()
    
}





