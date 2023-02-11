package ui

type Stack struct {
	Push   func(window *Window)
	Pop    func() *Window
	Peek   func() *Window
	Length func() int
}

func NewStack() Stack {
	slice := make([]*Window, 0)
	return Stack{
		Push: func(i *Window) {
			slice = append(slice, i)
		},
		Pop: func() *Window {
			res := slice[len(slice)-1]
			slice = slice[:len(slice)-1]
			return res
		},
		Peek: func() *Window {
			if len(slice)-1 < 0 {
				return nil
			}
			return slice[len(slice)-1]
		},
		Length: func() int {
			return len(slice)
		},
	}
}
