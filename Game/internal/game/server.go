package game

import (
	"sync"
)

type Server struct{
	Rooms sync.Map
	Users sync.Map
}

func CreateServer() *Server{
	return &Server{}
}

func (s *Server) AddRoom(r *Room){
	s.Rooms.Store(r.ID, r)
}

func (s *Server) DeleteRoom(r_id string){
	s.Rooms.Delete(r_id)
}

func (s *Server) GetRoom(r_id string) *Room{
	r, ok := s.Rooms.Load(r_id)
	if !ok{
		return nil
	}
	return r.(*Room)
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

