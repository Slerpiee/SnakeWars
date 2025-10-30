package game

import (
	"sync"
	"time"
	"log"
	"fmt"
    "github.com/google/uuid"
)

type RoomState struct{
    IsFull bool
    Started bool
    Waiting bool
}

type RoomStats struct{
    PlayerCount int
    PlayersReady int
    MaxPlayers int
}

type RoomMessage struct {
    StatusCode int
    Text string
}

func CreateMessage(code int, text string) RoomMessage{
    return RoomMessage{code, text}
}

type Room struct{
	ID string //Room ID
    Name string


	Users sync.Map  //ID: *User

	room_mutex sync.RWMutex
	State RoomState
	Stats RoomStats

	roomInput chan RoomMessage
	ticker *time.Ticker //update per second 

}

func CreateRoom(name string) *Room{
    return &Room{
        ID: uuid.NewString(),
        Name: name,
        Users: sync.Map{},
        room_mutex: sync.RWMutex{},
        State: RoomState{},
        Stats: RoomStats{MaxPlayers: 2},
        roomInput: make(chan RoomMessage),
        ticker: nil,
    }
}

func (room *Room) CanStart()bool{
    room.room_mutex.Lock()
    defer room.room_mutex.Unlock()
    return !room.State.Started && room.Stats.PlayerCount == room.Stats.PlayersReady && room.Stats.PlayerCount == room.Stats.MaxPlayers
}



func (room *Room) addUser(user *User) {
    room.Users.Store(user.ID, user)
}

func (room *Room) removeUser(id string){
    room.Users.Delete(id)
}

func (room *Room) userJoin(user *User){
    room.room_mutex.Lock()
    room.Stats.PlayerCount++
    room.room_mutex.Unlock()
    room.addUser(user)
    room.Broadcast(CreateMessage(1, user.ID))
}

func (room *Room) userDisconnect(id string){
    room.room_mutex.Lock()
    room.Stats.PlayerCount--
    room.room_mutex.Unlock()
    cand := room.getUser(id)
    if cand != nil{
        room.removeUser(id)
    }
    room.Broadcast(CreateMessage(-1, id))
}


func (room *Room) getUser(id string) *User {
	value, ok := room.Users.Load(id)
    if ok {
        return value.(*User)
    }
    return nil
}


func (room *Room) UpdateUser(id string) {
    if User := room.getUser(id); User != nil {
        User.mutex.Lock() 
        defer User.mutex.Unlock()
		User.Snake.Move(time.Since(User.LastPing), false)
    }
}


func (room *Room) Broadcast(message interface{}) {
    room.Users.Range(func(key, value interface{}) bool {
        User := value.(*User)
        
        select {
        	case User.SendChan <- message:
        default:
            log.Printf("User %s is lagging", User.ID) //Пользователь не забирает сообщения из канала
        }
        return true
    })
}


func (room *Room) startInputProcessor() {
    room.roomInput = make(chan RoomMessage, 100)
    
    go func() {
        for message := range room.roomInput { //Апдейты типо пользователь нажал кнопку управления, кто-то ебнулся об стенку и сдох и тд
            //handle room update
			log.Printf("Message for room %s: ", room.ID)
			fmt.Println(message)
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



func (room *Room) StartRoom(){
    room.startInputProcessor()
    room.startGameLoop()
}


func (room *Room) updateSnakes(){
    room.Users.Range(func(_, value any) bool {
        User := value.(*User)
        snake := User.Snake
        snake.Move(time.Since(User.LastPing), false)
        return true //если функция возвращает false, то процесс прирывается
    })
}




// func (room *Room) GameCollissions() []map[string][[]string]{ 

// }

func (room *Room) Close(){
    room.room_mutex.Lock()
    defer room.room_mutex.Unlock()

    if room.ticker != nil{
        room.ticker.Stop()
    }

    if room.roomInput != nil{
        close(room.roomInput)
    }

    var wg sync.WaitGroup
    room.Users.Range(func(key, val any) bool {
        wg.Add(1)
        go func (u *User)  {
            defer wg.Done()
            if u != nil{
                u.Close()
            }
        }(val.(*User))
        return true
    })

    wg.Wait()
    room.Users = sync.Map{}
    log.Printf("Room %s successfully closed", room.ID)
}


func (room *Room) gameTick() {
    room.room_mutex.Lock()
    defer room.room_mutex.Unlock()
	room.updateSnakes()
	//if checkCollisions -> kill
	//... BroadCast informatoin about kill or sum
    
}





