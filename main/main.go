package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/fzipp/canvas"
)

type Node struct {
	row int
	col int
}

type Grid struct{
	Squares []Square
	Ctx     *canvas.Context
	X,Y     float64
}

type Square struct {
	Side_front, Side_back, Side_left, Side_right bool
	Node bool
	Id, Row, Col                                 int
	Ctx                                          *canvas.Context
}

//CREATE A LABIRINT
const rows int = 10
const cols int = 10
const sides int = 4
const offset float64 = 20.0


func main() {
	http := flag.String("http", ":8080", "HTTP service address (e.g., '127.0.0.1:8080' or just ':8080')")
	flag.Parse()
	fmt.Println("Listening on " + httpLink(*http))
	err := canvas.ListenAndServe(*http, run,
		canvas.Title("Labyrinth"),
		canvas.Size(500, 500),
		canvas.ScaleFullPage(false, true),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// run function
func run(ctx *canvas.Context) {
	var h = &Grid{Ctx: ctx}
	rand.Seed(time.Now().UTC().UnixNano())
	lab := CreateStarter()
	ctx.SetLineWidth(2)
	ctx.SetStrokeStyleString("#8806ce")
	h.Squares = lab
	fmt.Println("Start animation...")
	for {
		select {
		case event := <-ctx.Events():
			if _, ok := event.(canvas.CloseEvent); ok {
				fmt.Println("Close Event...")
				return
			}
		default:
			if exit(lab){
				fmt.Println("Succefully finished the algoritm...")
				h.AddStartStop(ctx)
				ctx.Flush()
				h.PrintGrid()
				h.GetNodes()
				h.PrintNodes()
				return
			}else{
				h.update()
				h.draw(ctx)
				ctx.Flush()
				time.Sleep(time.Second / 1500)
			}
		}
	}
}


// update func
func (h *Grid) update(){
	h.Squares = RandomSquares(h.Squares)
}

// create a starter Maze
func CreateStarter()[]Square {
	var id int = 1
	grid := make([]Square,rows*cols)
	for row:=0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			for side := 0; side < sides; side++ {
				switch side {
				case 0:
					grid[id-1].Side_front = true
				case 1:
					grid[id-1].Side_right = true
				case 2:
					grid[id-1].Side_left = true
				case 3:
					grid[id-1].Side_back = true
				}
			}
			grid[id-1].Id = id
			grid[id-1].Col = col + 1
			grid[id-1].Row = row + 1
			id++
		}
	}
	return grid
}

// draw labyrinth
func (h *Grid) draw(ctx *canvas.Context) {
	ctx.ClearRect(0,0,1000,1000)
	h.X = 50.0
	h.Y = 50.0
	var count int
	for row:=1; row <= rows; row++ {
		for col := 1; col <= cols; col++ {
			for side := 0; side < sides; side++ {
				switch side {
				case 0:
					// back side
					if h.Squares[count].Side_back {
						h.Ctx.BeginPath()
						h.Ctx.MoveTo(h.X, h.Y)
						h.Ctx.LineTo(h.X+offset, h.Y)
						h.Ctx.Stroke()
						h.Ctx.ClosePath()
						h.X = h.X + offset
					}else{
						h.X = h.X + offset
					}
				case 1:
					// side right
					if h.Squares[count].Side_right {
						h.Ctx.BeginPath()
						h.Ctx.MoveTo(h.X, h.Y)
						h.Ctx.LineTo(h.X, h.Y-offset)
						h.Ctx.Stroke()
						h.Ctx.ClosePath()
						h.Y = h.Y - offset
					}else{
						h.Y = h.Y - offset
					}
				case 2:
					// side front
					if h.Squares[count].Side_front {
						h.Ctx.BeginPath()
						h.Ctx.MoveTo(h.X, h.Y)
						h.Ctx.LineTo(h.X-offset, h.Y)
						h.Ctx.Stroke()
						h.Ctx.ClosePath()
						h.X = h.X - 20
					}else{
						h.X = h.X - 20
					}
				case 3:
					// side left
					if h.Squares[count].Side_left {
						h.Ctx.BeginPath()
						h.Ctx.MoveTo(h.X, h.Y)
						h.Ctx.LineTo(h.X, h.Y+offset)
						h.Ctx.Stroke()
						h.Ctx.ClosePath()
						h.Y = h.Y + offset
					}else{
						h.Y = h.Y + offset
					}
				}
			}
			h.X = h.X + offset
			count++
		}
		h.X = h.X - float64(cols)*offset
		h.Y = h.Y + offset
	}
}

// exit condition
func exit(grid []Square)bool{
	var target int
	for i, square := range grid{
		if i==0{
			target = square.Id
		}else{
			if square.Id !=target{
				return false
			}
		}
	}
	return  true
}

// generate 2 random num that are Id's adjoins
func RandomSquares(grid []Square)[]Square {
	for{
		random_row1:=rand.Intn(rows)+1
		random_row2:=rand.Intn(rows)+1
		random_col1:=rand.Intn(cols)+1
		random_col2:=rand.Intn(cols)+1

		//check if random same cell
		if random_row1==random_row2 && random_col1==random_col2{
			continue
		}

		adjoins := Adjoins(random_row1,random_row2,random_col1,random_col2)
		differentId,id1,id2 := DifferentId(random_row1,random_row2,random_col1,random_col2,grid)
			if adjoins && differentId{
				BreakWall(grid,id1,id2,random_row1,random_row2,random_col1,random_col2)
				break
			}
		}
	return grid
}

// check if 2 cell are adjoins
func Adjoins(row1,row2,col1,col2 int)bool {
	return SameCol(col1, col2, row1, row2) || SameRow(col1, col2, row1, row2)
}

// check if cell are in the same row
func SameRow(col1,col2,row1,row2 int)bool{
	return row1==row2 && math.Abs(float64(col2-col1))==1
}

// check if cell are in the same Col
func SameCol(col1,col2,row1,row2 int)bool{
	return col1==col2 && math.Abs(float64(row2-row1))==1
}

// check if Id are different
func DifferentId(r1,r2,c1,c2 int,grid []Square)(bool,int,int){
	var id1,id2 int
	for _,square := range grid{
		if square.Row == r1 && square.Col == c1{
			id1 = square.Id
		}
		if square.Row == r2 && square.Col == c2{
			id2 = square.Id
		}
	}
	return id1!=id2,id1,id2
}

// break the "wall" in the middle of 2 cell
func BreakWall(grid []Square, num1,num2,row1,row2,col1,col2 int){
	var max int = Max(num1,num2)
	var min int = Min(num1,num2)

	for i, square := range grid{
		if square.Row == row1 && square.Col == col1 && square.Id == max{
			switch {
			case row1>row2:
				grid[i].Side_front = false
				grid[i-cols].Side_back = false
			case col1>col2:
				grid[i].Side_left = false
				grid[i-1].Side_right = false
			case row2>row1:
				grid[i].Side_back = false
				grid[i+cols].Side_front = false
			case col2>col1:
				grid[i].Side_right = false
				grid[i+1].Side_left = false
			}
			grid = ChangeValues(max,min,grid)
		}
	}
}

// change all max to min
func ChangeValues(max int,min int,grid []Square)[]Square {
	for i, square:= range grid{
		if square.Id == max{
			grid[i].Id = min
		}
	}
	return grid
}

// max of 2 num
func Max(num1,num2 int)int{
	if num1>num2{
		return num1
	}
	return num2
}

//min of 2 num
func Min(num1,num2 int)int{
	if num1<num2{
		return num1
	}
	return num2
}

// print the Maze data in terminal
func (h *Grid)PrintGrid(){
	for _,square := range h.Squares{
		fmt.Println(square)
	}
}

// link for localhost
func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}

func (h *Grid)AddStartStop(ctx *canvas.Context){
	h.Ctx.SetFont("10px Arial")
	h.Ctx.FillText("St", 52, 45)
	h.Ctx.FillText("Fi", 235, 225)
}

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


func (h *Grid)GetNodes(){
	direction := map[string]bool{"up":false, "back":false, "right":false, "left":false}

	//first/last node
	h.Squares[0].Node = true
	h.Squares[len(h.Squares)-1].Node = true

	//direction of the start
	if h.Squares[0].Side_right==false {
		direction["right"] = true
	}
	if h.Squares[0].Side_back==false{
		direction["back"] = true
	}
	//other nodes
	fmt.Println("Starting direction:",direction)
	for i:=1;i<len(h.Squares);i++{
		if direction["right"] {
			h.GoRight()
			direction["right"] = false
		}else if direction["back"]{
			h.GoBack()
			direction["back"] = false
		}
	}
	return
}

func (h *Grid)GoBack(){
	fmt.Print("Back: ")
	for i:=cols;i<rows*cols;i+=cols{
		if h.Squares[i].Side_back{
			h.Squares[i].Node = true
			fmt.Print("(",h.Squares[i].Row,",",h.Squares[i].Col,")")
			break
		}
	}
	fmt.Println()
}

func (h *Grid)GoRight(){
	fmt.Print("Right: ")
	for i:=0;i<cols;i++{
		if h.Squares[i].Side_right{
			h.Squares[i].Node = true
			fmt.Print("(",h.Squares[i].Row,",",h.Squares[i].Col,")")
			break
		}
	}
	fmt.Println()
}

func (h *Grid)PrintNodes(){
	fmt.Print("Nodes: ")
	for i:=0;i<len(h.Squares);i++{
		if h.Squares[i].Node{
			fmt.Print("(",h.Squares[i].Row,",",h.Squares[i].Col,") ")
		}
	}
}