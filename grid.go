package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fzipp/canvas"
)

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

func run(ctx *canvas.Context) {
	ctx.SetLineWidth(2)
	ctx.SetStrokeStyleString("#8806ce")
	h := &labyrinth{ctx: ctx}
	h.draw()
	ctx.Flush()
}

type labyrinth struct {
	ctx  *canvas.Context
	x, y float64
}

func (h *labyrinth) draw() {
	h.x = 50.0
	h.y = 50.0
	for k:=0;k<20;k++{
		for i := 0; i < 20; i++ {
			for j := 0; j < 4; j++ {
				switch j {
				case 0:
					h.ctx.BeginPath()
					h.ctx.MoveTo(h.x, h.y)
					h.ctx.LineTo(h.x+20, h.y)
					h.ctx.Stroke()
					h.ctx.ClosePath()
					h.x = h.x +20
				case 1:
					h.ctx.BeginPath()
					h.ctx.MoveTo(h.x, h.y)
					h.ctx.LineTo(h.x, h.y-20)
					h.ctx.Stroke()
					h.ctx.ClosePath()
					h.y = h.y -20
				case 2:
					h.ctx.BeginPath()
					h.ctx.MoveTo(h.x, h.y)
					h.ctx.LineTo(h.x-20, h.y)
					h.ctx.Stroke()
					h.ctx.ClosePath()
					h.x = h.x - 20
				case 3:
					h.ctx.BeginPath()
					h.ctx.MoveTo(h.x, h.y)
					h.ctx.LineTo(h.x, h.y+20)
					h.ctx.Stroke()
					h.ctx.ClosePath()
					h.y = h.y +20
				}
			}
			h.x = h.x + 20
		}
		h.x = h.x - 20*20
		h.y = h.y +20
	}


}

func httpLink(addr string) string {
	if addr[0] == ':' {
		addr = "localhost" + addr
	}
	return "http://" + addr
}
