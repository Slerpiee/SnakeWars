package game

import (
	"sync"
)

type Server struct{
	Rooms sync.Map
	Users sync.Map
	Names map[string]string //ID:ROOM_ID
	mutex sync.RWMutex
}

func CreateServer() *Server{
	return &Server{Names: make(map[string]string)}
}

func (s *Server) AddRoom(r *Room){
	s.Rooms.Store(r.ID, r)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Names[r.Name] = s.Names[r.ID]
}

func (s *Server) DeleteRoom(r_id string){
	temp := s.GetRoom(r_id)

	if temp != nil{
		s.mutex.Lock()
		delete(s.Names, temp.Name)
		s.mutex.Unlock()
	}

	s.Rooms.Delete(r_id)
	
}

func (s *Server) GetRoom(r_id string) *Room{
	r, ok := s.Rooms.Load(r_id)
	if !ok{
		return nil
	}
	return r.(*Room)
}

func(s *Server) GetRoomIdByName(name string) (string, bool){
	s.mutex.RLock()
	cand, ok := s.Names[name]
	s.mutex.RUnlock()
	return cand, ok
}

func (s *Server) GetRooms() map[string]string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	

	result := make(map[string]string, len(s.Names))

	for k, v := range s.Names {
		result[k] = v
	}
	return result
}


func (s *Server) AddUser(u *User){
	s.Users.Store(u.ID, u)
}
func (s *Server) DeleteUser(u_id string){
	s.Users.Delete(u_id)
}

func (s *Server) GetUser(u_id string) *User{
	u, ok := s.Users.Load(u_id)
	if !ok{
		return nil
	}
	return u.(*User)
}

