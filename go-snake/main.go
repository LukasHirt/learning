package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
)

type SnakePart struct {
	X int
	Y int
}

type Snake struct {
	Parts []SnakePart
	XSpeed int
	YSpeed int
}

type Game struct {
	screen tcell.Screen
	snake *Snake
}

func ChangeDirection(snake *Snake, x int, y int) {
 	snake.XSpeed = x
	snake.YSpeed = y
}

func Move(snake *Snake) {
	snake.Parts = snake.Parts[1:]

	if snake.XSpeed != 0 {
		snake.Parts = append(snake.Parts, SnakePart { snake.Parts[len(snake.Parts) - 1].X + snake.XSpeed, snake.Parts[len(snake.Parts) - 1].Y })
		return
	}

	if snake.YSpeed != 0 {
		snake.Parts = append(snake.Parts, SnakePart { snake.Parts[len(snake.Parts) - 1].X, snake.Parts[len(snake.Parts) - 1].Y + snake.YSpeed })
		return
	}
}

func DrawParts(screen tcell.Screen, parts []SnakePart) {
	box_style := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)

	for _, part := range parts {
		screen.SetContent(part.X, part.Y, ' ', nil, box_style)
	}
}

func Run(game *Game) {
	game_style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	game.screen.SetStyle(game_style)

	for {
		game.screen.Clear()

		Move(game.snake)
		DrawParts(game.screen, game.snake.Parts)

		time.Sleep(40 * time.Millisecond)
		game.screen.Show()
	}
}

func main() {
	screen, err := tcell.NewScreen()

	if err != nil {
		fmt.Println("Error creating screen: ", err)
		panic(err)
	}

	if err := screen.Init(); err != nil {
		fmt.Println("Error initializing screen: ", err)
		panic(err)
	}

	quit := func() {
		maybe_panic := recover()
		screen.Fini()

		if (maybe_panic != nil) {
			panic(maybe_panic)
		}
	}

	defer quit()

	snake_parts := []SnakePart { { 4, 4 }, { 5, 4 }, { 6, 4 } }
	snake := Snake { snake_parts, 1, 0 }
	game := Game { screen, &snake }

	go Run(&game)

	for {
		ev := screen.PollEvent()

		switch ev := ev.(type) {
			case *tcell.EventResize: {
				screen.Sync()
			}

			case *tcell.EventKey: {
				if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEscape {
					return
				}

				if ev.Key() == tcell.KeyUp {
					ChangeDirection(&snake, 0, -1)
				}

				if ev.Key() == tcell.KeyDown {
					ChangeDirection(&snake, 0, 1)
				}

				if ev.Key() == tcell.KeyLeft {
					ChangeDirection(&snake, -1, 0)
				}

				if ev.Key() == tcell.KeyRight {
					ChangeDirection(&snake, 1, 0)
				}
			}
		}
	}
}