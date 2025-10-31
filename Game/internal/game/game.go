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
)

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
	ID string  `json:"id"`
    Name string `json:"name"`


	Users sync.Map  `json:"-"` //USER_ID: *User

	room_mutex sync.RWMutex `json:"-"`
	State RoomState 
	Stats RoomStats `json:"stats"`

	roomInput chan RoomMessage `json:"-"`
	ticker *time.Ticker `json:"-"`

}

func CreateRoom(name string) *Room{
    return &Room{
        ID: uuid.NewString(),
        Name: name,
        Users: sync.Map{},
        room_mutex: sync.RWMutex{},
        State: ROOMSTATE_WAITING,
        Stats: RoomStats{},
        roomInput: make(chan RoomMessage),
        ticker: nil,
    }
}

func (room *Room) CanStart()bool{
    room.room_mutex.Lock()
    defer room.room_mutex.Unlock()
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
    room.room_mutex.Lock()
    room.State = ROOMSTATE_GAME_STARTED
    room.room_mutex.Unlock()

    room.InitSnakes()
    room.Broadcast(CreateMessage(ROOMSTATE_GAME_STARTED, "Game Start"))
    log.Printf("Game started in room %s", room.ID)
}

func (room *Room) EndGame(){
    room.room_mutex.Lock()
    room.State = ROOMSTATE_GAME_ENDED
    room.room_mutex.Unlock()
    
    room.Broadcast(CreateMessage(ROOMSTATE_GAME_ENDED, "GAME_END"))
}

func (room *Room) InitSnakes(){
    positions := []Point{
        {X: 100, Y: 100},
        {X: 900, Y: 900}, //пока для 2 пользователей, надо будет для большего кол-ва дописать генерацию
    }
    i := 0
    room.Users.Range(func(key, value interface{}) bool{
        if i < len(positions){
            user := value.(*User)
            user.mutex.Lock()
            user.Snake.Head = positions[i]
            user.Snake.State.isAlive = true
            user.mutex.Unlock()
            i++
        }
        return true
    })
}

func (room *Room) AddUser(user *User) {
    room.Users.Store(user.ID, user)
}

func (room *Room) RemoveUser(id string){
    room.Users.Delete(id)
}

func (room *Room) UserJoin(user *User){
    room.room_mutex.Lock()
    room.Stats.PlayerCount++ //Можно попробовать без лока комнаты увеличить атомарно
    room.room_mutex.Unlock()
    room.AddUser(user)
    room.Broadcast(CreateMessage(1, user.ID))
}

func (room *Room) UserDisconnect(id string){
    room.room_mutex.Lock()
    room.Stats.PlayerCount--
    //room.State.IsFull = false
    if room.Stats.PlayerCount < 2 && room.State == ROOMSTATE_GAME_STARTED{
        room.State = ROOMSTATE_GAME_ENDED
        room.Broadcast(CreateMessage(10, "NOT ENOUGH PLAYERS"))
    }
    room.room_mutex.Unlock()
    cand := room.GetUser(id)
    if cand != nil{
        room.RemoveUser(id)
    }
    room.Broadcast(CreateMessage(-1, id))
}


func (room *Room) GetUser(id string) *User {
	value, ok := room.Users.Load(id)
    if ok {
        return value.(*User)
    }
    return nil
}


func (room *Room) UpdateUser(id string) {
    if User := room.GetUser(id); User != nil {
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
            room.gameTick() //Добавить контекст отмены или пока хзе
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
    room.updateSnakes()
    //
}





