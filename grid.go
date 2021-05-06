package main

import (
	"flag"
	"fmt"
	"github.com/fzipp/canvas"
	"log"
	"math"
	"math/rand"
)

/*
PART 1 -CREATE A LABYRINTH-
Matrix nxm associate each cell a num
Choose randomly 2 cell that has different value and that are adjoins
Change the value of the max of the 2 num of the one lower and all the other max values
Also brak the wall between them
Continue until all numbers are equal

PART 2-SOLVE A LABYRINTH-
Dijkstra's algoritm to resolve it
 */

const rows int = 20
const cols int = 20
const sides int = 4
const offset float64 = 20.0

//main func
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

//run function
func run(ctx *canvas.Context) {
	ctx.SetLineWidth(2)
	ctx.SetStrokeStyleString("#8806ce")
	h := &labyrinth{ctx: ctx}
	var grid []square
	h.draw(grid)
	ctx.Flush()
}

//struct for cell
type square struct{
	side_front,side_back,side_left,side_right bool
	id,row,col int
}

//struct for cells
type labyrinth struct {
	ctx  *canvas.Context
	x, y float64
}

//draw labyrinth
func (h *labyrinth) draw(grid []square) {
	h.x = 50.0
	h.y = 50.0
	var id int = 1
	grid = make([]square,rows*cols)
	for row:=0; row < rows; row++{
		for col := 0; col < cols; col++ {
			for side := 0; side < sides; side++ {
				switch side {
				case 0:
					h.ctx.BeginPath()
					h.ctx.MoveTo(h.x, h.y)
					h.ctx.LineTo(h.x + offset, h.y)
					h.ctx.Stroke()
					h.ctx.ClosePath()
					h.x = h.x + offset
					grid[id-1].side_front = true
				case 1:
					h.ctx.BeginPath()
					h.ctx.MoveTo(h.x, h.y)
					h.ctx.LineTo(h.x, h.y - offset)
					h.ctx.Stroke()
					h.ctx.ClosePath()
					h.y = h.y - offset
					grid[id-1].side_right = true
				case 2:
					h.ctx.BeginPath()
					h.ctx.MoveTo(h.x, h.y)
					h.ctx.LineTo(h.x - offset, h.y)
					h.ctx.Stroke()
					h.ctx.ClosePath()
					h.x = h.x - 20
					grid[id-1].side_left = true
				case 3:
					h.ctx.BeginPath()
					h.ctx.MoveTo(h.x, h.y)
					h.ctx.LineTo(h.x, h.y + offset)
					h.ctx.Stroke()
					h.ctx.ClosePath()
					h.y = h.y + offset
					grid[id-1].side_back = true
				}
			}
			grid[id-1].id = id
			grid[id-1].col = col + 1
			grid[id-1].row = row + 1
			id++
			h.x = h.x + offset
		}
		h.x = h.x - float64(cols) * offset
		h.y = h.y + offset
	}
	PrintGrid(grid)
}

//generate 2 random num that are id's adjoins
func RandomSquares(grid []square){
	for{
		random_num1:=rand.Intn(rows*cols)+1
		random_num2:=rand.Intn(rows*cols)+1
		if ok,row1,col1,row2,col2:=SearchRandomInGrid(random_num1,random_num2,grid);ok {
			if Adjoins(row1,col1,row2,col2) && DifferentId(random_num1,random_num2) {
				BreakWall(grid,random_num1,random_num2,row1,row2,col1,col2)
				break
			}
		}
	}
	PrintGrid(grid)
}

//search if 2 num random are both in id of grid
func SearchRandomInGrid(num1,num2 int, grid []square)(bool,int,int,int,int){
	var target_one,target_two bool
	var target_one_row,target_one_col int
	var target_two_row,target_two_col int

	for _,square:= range grid{
		if square.id == num1{
			target_one = true
			target_one_row = square.row
			target_one_col = square.col
		}
		if square.id == num2{
			target_two = true
			target_two_row = square.row
			target_two_col = square.col
		}
	}
	return target_one && target_two,target_one_row,target_one_col,target_two_row,target_two_col
}

//check if 2 cell are adjoins
func Adjoins(row1,row2,col1,col2 int)bool {
	return SameCol(col1, col2, row1, row2) || SameRow(col1, col2, row1, row2)
}

//check if cell are in the same row
func SameRow(col1,col2,row1,row2 int)bool{
	return row1==row2 && math.Abs(float64(col2-col1))==1
}

//check if cell are in the same col
func SameCol(col1,col2,row1,row2 int)bool{
	return col1==col2 && && math.Abs(float64(row2-row1))==1
}

//check if id are different
func DifferentId(id1,id2 int)bool{
	return id1!=id2
}

//break the "wall" in the middle of 2 cell
//not finished
func BreakWall(grid []square, num1,num2,row1,row2,col1,col2 int){
	var max int = Max(num1,num2)
	var min int = Min(num1,num2)

	for i, square := range grid{
		if square.row == row1 && square.col == col1 && square.id == max{
			switch {
			case row1>row2:
				grid[i].side_front = false
				grid[i-cols].side_back = false
			case col1>col2:
				grid[i].side_left = false
				grid[i-1].side_right = false
			case row2>row1:
				grid[i].side_back = false
				grid[i+cols].side_front = false
			case col2>col1:
				grid[i].side_right = false
				grid[i+1].side_left = false
			}
		}
		if square.id == max{
			grid[i].id = min
		}
	}
}

func Max(num1,num2 int)int{
	if num1>num2{
		return num1
	}
	return num2
}

func Min(num1,num2 int)int{
	if num1<num2{
		return num1
	}
	return num2
}

//Print the grid data in terminal
func PrintGrid(grid []square){
	for _,square := range grid{
		fmt.Println(square)
	}
}

//link for localhost
func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}
