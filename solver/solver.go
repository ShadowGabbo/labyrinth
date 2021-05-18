package solver

import "github.com/fzipp/canvas"

//Breadth-first search algorithm
	//getting nodes:
	//	-start from the start
	//	-look down/sx/up/dx
	//	-boolean array if visited
	//	-for rows:
	//		-first node at the start
	//		-if im in a path and i was in a wall is a node
	//		-if i can go up or down is a node
	//		-if im in a path and i will be in a wall is a node
	//		-if i cant go left or right isnt a node, and so on...
	//		-count node the min is the solver


type Grid struct{
	Squares []Square
	Ctx     *canvas.Context
	X,Y     float64
}

type Square struct {
	Side_front, Side_back, Side_left, Side_right bool
	Id, Row, Col                                 int
	Ctx                                          *canvas.Context
}

//func GetNodes(label *Grid){
//}