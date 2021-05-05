package main

import (
	"fmt"
	"github.com/fzipp/canvas"
)

const rows int = 10
const cols int = 20

func main() {
	//struttura dati n x m
	//ogni cella divide un muro con le 4 adiacenti tranne il bordo 3 e nell angolo 2
	//dentro ogni cella cè un numero intero diverso per cella
	//scelgo 2 celle adiacenti che abbiano numeri fra di loro diversi e butto giu il muro
	//che le divide, facendo diventare tutti i numeri sulla scacchiera uguali  al più alto
	//dei due, uguali all'altro
	//ripeto l'operazione finche non mi rimane solo un numero in tutte le celle
	StarterTemplate := Starter(rows, cols)
	InizializedTemplate := Inizialize(StarterTemplate)
	PrintStarter(InizializedTemplate)
}

//make a starter template given cols/rows
func Starter(rows int, cols int) [][]int {
	labyrinth := make([][]int, rows)
	for i := 0; i < rows; i++ {
		labyrinth[i] = make([]int, cols)
	}
	return labyrinth
}

//print the template
func PrintStarter(labyrinth [][]int) {
	for index_row, row := range labyrinth {

		//print first line
		for i := 0; i < 4*cols; i++ {
			fmt.Print("-")
		}
		fmt.Println()

		//print the body
		for _, col := range row {
			fmt.Print("|", col)
			if ((col-100)+1)%cols == 0 {
				fmt.Print("|")
			}
		}
		fmt.Println()

		//print last line
		if index_row == len(labyrinth)-1 {
			for i := 0; i < 4*cols; i++ {
				fmt.Print("-")
			}
			fmt.Println()
		}
	}
	ctx.rect(x, y, width, height)
}

//inizializing the template
func Inizialize(labyrinth [][]int) [][]int {
	var counter int = 100
	for row := 0; row < len(labyrinth); row++ {
		for col := 0; col < len(labyrinth[row]); col++ {
			labyrinth[row][col] = counter
			counter++
		}
	}
	return labyrinth
}
