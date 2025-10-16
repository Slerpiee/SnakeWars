package game
const SNAKE_WIDTH = 5

type Players struct{
	player1 *Snake
	player2 *Snake
}

func CheckHeadsCollision(s1, s2 *Snake) bool{
	if s1.ID == s2.ID{
		return CheckSelfCollision(s1)
	}
	head1 := s1.Head
	head2 := s2.Head
	dist_heads := distance(head1, head2)
	if dist_heads <= SNAKE_WIDTH*2{
		return true
	}
	return false
}

func CheckSelfCollision(s1 *Snake) bool{
	head := s1.Head
	size := s1.Body.Size()
	body := s1.Body
	if size <= 1{
		return false
	}
	for i := size - 1; i >= 0; i --{
		cur_seg := body[i] //переделать
		dist := distancePointLine(head, cur_seg)
		
	}
}

func CheckBodyCollision(s1, s2 *Snake) Players{

}