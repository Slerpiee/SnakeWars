package game

import (
	"sync"
	"time"
	"log"
	"fmt"
    "github.com/google/uuid"
)

type RoomState int



const (
    ROOMSTATE_WAITING = 0
    ROOMSTATE_GAME_STARTED = 1
    ROOMSTATE_GAME_ENDED = 2  
    ROOM_PLAYER_CONNECTED = 3
    ROOM_PLAYER_DISCONNECTED=4 
    ROOM_PLAYER_UPDATE = 5
    ROOMSTATE_ROOM_CLOSED = 6
)

type RoomStats struct{
    PlayerCount int
    PlayersReady int
    MaxPlayers int
}





type Room struct{
	ID string  `json:"id"`
    Name string `json:"name"`
	Users map[string]*User  `json:"-"` //USER_ID: *User
	mutex sync.RWMutex `json:"-"`
	State RoomState 
	Stats RoomStats `json:"stats"`
	roomInput chan Message `json:"-"`
	ticker *time.Ticker `json:"-"`

}

func CreateRoom(name string) *Room{
    return &Room{
        ID: uuid.NewString(),
        Name: name,
        Users: map[string]*User{},
        mutex: sync.RWMutex{},
        State: ROOMSTATE_WAITING,
        Stats: RoomStats{},
        roomInput: make(chan Message),
        ticker: nil,
    }
}


func (room *Room) CanStart()bool{
    room.mutex.RLock()
    defer room.mutex.RUnlock()
    return !(room.State == ROOMSTATE_WAITING)&& room.Stats.PlayerCount == room.Stats.PlayersReady && room.Stats.PlayerCount == room.Stats.MaxPlayers
}


func (room *Room) PlayerSetReady(id string, ready bool){
    if user := room.GetUser(id); user != nil{
        user.mutex.Lock()
        user.Snake.State.isReady = ready
        user.mutex.Unlock()
        //room.UpdateReady()
    }
}


func (room *Room) StartGame(){
    room.mutex.Lock()
    room.State = ROOMSTATE_GAME_STARTED
    room.mutex.Unlock()

    room.InitSnakes()
    room.Broadcast(CreateMessage(room.ID, ROOMSTATE_GAME_STARTED, nil))
    log.Printf("Game started in room %s", room.ID)
}


func (room *Room) EndGame(){
    room.mutex.Lock()
    room.State = ROOMSTATE_GAME_ENDED
    room.mutex.Unlock()
    room.Broadcast(CreateMessage(room.ID, ROOMSTATE_GAME_ENDED, nil))
}


func (room *Room) InitSnakes(){
    positions := []Point{
        {X: 100, Y: 100},
        {X: 900, Y: 900}, //пока для 2 пользователей, надо будет для большего кол-ва дописать генерацию
    }
    i := 0
    room.mutex.Lock()
    defer room.mutex.Unlock()
    for _, user := range room.Users{
        if i < len(positions){
            user.mutex.Lock()
            user.Snake.Head = positions[i] 
            user.Snake.State.isAlive = true
            user.mutex.Unlock()
            i++
        }
    }
}


func (room *Room) AddUser(user *User) bool{ //existed before?
    room.mutex.Lock()
    defer room.mutex.Unlock()
    _, ok := room.Users[user.ID]
    if !ok{
        room.Users[user.ID] = user
    }
    return ok
}

func (room *Room) RemoveUser(id string){ 
    room.mutex.Lock()
    defer room.mutex.Unlock()
    delete(room.Users, id)
}


func (room *Room) UserJoin(user *User){
    room.mutex.Lock()
    room.Stats.PlayerCount++ //Можно попробовать без лока комнаты увеличить атомарно
    room.mutex.Unlock()
    room.AddUser(user)
    type m struct{
        user_id string 
    }
    message := m{user_id: user.ID}
    room.Broadcast(CreateMessage(room.ID, ROOM_PLAYER_CONNECTED, message))
}

func (room *Room) UserDisconnect(id string){
    room.mutex.Lock()
    room.Stats.PlayerCount--
    room.mutex.Unlock()
    room.EndGame()  
    room.RemoveUser(id)
    room.Broadcast(CreateMessage(room.ID, ROOM_PLAYER_DISCONNECTED, id))
    
}


func (room *Room) GetUser(id string) *User {
	room.mutex.RLock()
    defer room.mutex.Unlock()
    cand, ok := room.Users[id]
    if !ok{
        return nil
    }
    return cand
}


func (room *Room) UpdateUser(id string) {
    if User := room.GetUser(id); User != nil {
        User.mutex.Lock() 
        defer User.mutex.Unlock()
		User.Snake.Move(time.Since(User.LastPing), false)
    }
}


func (room *Room) Broadcast(message Message) {
    room.mutex.RLock()
    defer room.mutex.RUnlock()
    for _, user := range room.Users {
        select {
        	case user.Chan <- message:
        default:
            log.Printf("User %s is lagging", user.ID) //Пользователь не забирает сообщения из канала
        }
    }
}


func (room *Room) startInputProcessor() {
    room.roomInput = make(chan Message, 100)


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
            room.gameTick() //Добавить контекст отмены или пока хзе
        }
    }()
}

func (room *Room) StartRoom(){
    room.startInputProcessor()
    room.startGameLoop()
}

func (room *Room) Close(){
    room.mutex.Lock()
    defer room.mutex.Unlock()

    if room.ticker != nil{
        room.ticker.Stop()
    }

    if room.roomInput != nil{
        close(room.roomInput)
    }

    room.Broadcast(CreateMessage(room.ID, ROOMSTATE_ROOM_CLOSED, nil))

    log.Printf("Room %s successfully closed", room.ID)
}



func (room *Room) updatePositions(){
    room.mutex.RLock()
    defer room.mutex.RUnlock()
    for _, user := range room.Users{
        user.UpdatePos()
    }
}


func (room *Room) gameTick() {
    room.updatePositions()
    //room.checkCollisions?
}





