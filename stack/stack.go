package stack

type Stack struct {
	slice []uint64
}

func (s *Stack) Push(b uint64) {
	s.slice = append(s.slice, b)
}

func (s *Stack) Pop() uint64 {
	v := s.Top()
	s.slice = s.slice[:len(s.slice)-1]
	return v
}

