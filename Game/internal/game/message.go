package game

type Message struct {
	ID string `json:"ID,omitempty"`
    Type    int     	`json:"type"`
    Room    string      `json:"room,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Sender  string      `json:"sender,omitempty"`
}

func CreateMessage(r string, t int, d interface{}) Message{
    return Message{
        Room: r,
        Type: t,
        Data: d,
    }
}