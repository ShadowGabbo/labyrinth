package main

import (
	"flag"
	"fmt"
	"github.com/fzipp/canvas"
	"log"
	"math"
	"math/rand"
	"time"
)

/*
PART 1 -CREATE A LABYRINTH-

PART 2-SOLVE A LABYRINTH-
Dijkstra's algoritm to resolve it
 */

const rows int = 10
const cols int = 10
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
	grid := CreateStarter()
	ctx.SetLineWidth(2)
	ctx.SetStrokeStyleString("#8806ce")
	h := &labyrinth{ctx: ctx}
	fmt.Println("Start animation...")
	for {
		select {
		case event := <-ctx.Events():
			if _, ok := event.(canvas.CloseEvent); ok {
				fmt.Println("Close Event...")
				return
			}
		default:
			if exit(grid){
				PrintGrid(grid)
				fmt.Println("Succefully finished the algoritm...")
				return
			}else{
				grid = RandomSquares(grid)
				h.draw(grid,ctx)
				ctx.Flush()
				time.Sleep(time.Second / 2)
			}
		}
	}
}

type square struct{
	side_front,side_back,side_left,side_right bool
	id,row,col int
}

type labyrinth struct {
	ctx  *canvas.Context
	x, y float64
}

func CreateStarter()[]square{
	var id int = 1
	grid := make([]square,rows*cols)
	for row:=0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			for side := 0; side < sides; side++ {
				switch side {
				case 0:
					grid[id-1].side_front = true
				case 1:
					grid[id-1].side_right = true
				case 2:
					grid[id-1].side_left = true
				case 3:
					grid[id-1].side_back = true
				}
			}
			grid[id-1].id = id
			grid[id-1].col = col + 1
			grid[id-1].row = row + 1
			id++
		}
	}
	return grid
}

//draw labyrinth
func (h *labyrinth) draw(grid []square,ctx *canvas.Context) {
	h.x = 50.0
	h.y = 50.0
	var count int
	for row:=1; row <= rows; row++ {
		for col := 1; col <= cols; col++ {
			for side := 0; side < sides; side++ {
				switch side {
				case 0:
					if grid[count].side_front {
						h.ctx.BeginPath()
						h.ctx.MoveTo(h.x, h.y)
						h.ctx.LineTo(h.x+offset, h.y)
						h.ctx.Stroke()
						h.ctx.ClosePath()
						h.x = h.x + offset
					}
				case 1:
					if grid[count].side_right {
						h.ctx.BeginPath()
						h.ctx.MoveTo(h.x, h.y)
						h.ctx.LineTo(h.x, h.y-offset)
						h.ctx.Stroke()
						h.ctx.ClosePath()
						h.y = h.y - offset
					}
				case 2:
					if grid[count].side_back {
						h.ctx.BeginPath()
						h.ctx.MoveTo(h.x, h.y)
						h.ctx.LineTo(h.x-offset, h.y)
						h.ctx.Stroke()
						h.ctx.ClosePath()
						h.x = h.x - 20
					}
				case 3:
					if grid[count].side_left {
						h.ctx.BeginPath()
						h.ctx.MoveTo(h.x, h.y)
						h.ctx.LineTo(h.x, h.y+offset)
						h.ctx.Stroke()
						h.ctx.ClosePath()
						h.y = h.y + offset
					}
				}
			}
			h.x = h.x + offset
			count++
		}
		h.x = h.x - float64(cols)*offset
		h.y = h.y + offset
	}
}

func exit(grid []square)bool{
	var target int
	for i, square := range grid{
		if i==0{
			target = square.id
		}else{
			if square.id!=target{
				return false
			}
		}
	}
	return  true
}

//generate 2 random num that are id's adjoins
func RandomSquares(grid []square)[]square{
	for{
		random_row1:=rand.Intn(rows)+1
		random_row2:=rand.Intn(rows)+1
		random_col1:=rand.Intn(cols)+1
		random_col2:=rand.Intn(cols)+1
		adjoins := Adjoins(random_row1,random_row2,random_col1,random_col2)
		differentId,id1,id2 := DifferentId(random_row1,random_row2,random_col1,random_col2,grid)
			if adjoins && differentId{
				BreakWall(grid,id1,id2,random_row1,random_row2,random_col1,random_col2)
				break
			}
		}
	return grid
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
	return col1==col2 && math.Abs(float64(row2-row1))==1
}

//check if id are different
func DifferentId(r1,r2,c1,c2 int,grid []square)(bool,int,int){
	var id1,id2 int
	for _,square := range grid{
		if square.row == r1 && square.col == c1{
			id1 = square.id
		}
		if square.row == r2 && square.col == c2{
			id2 = square.id
		}
	}
	return id1!=id2,id1,id2
}

//break the "wall" in the middle of 2 cell
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
