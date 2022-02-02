package main

Timport (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	flag.Parse()
	// initialize game
	initialise()
	defer cleanup()

	err := loadResource()
	if err != nil {
		log.Println("failed to load maze:", err)
		return
	}
	// load resources

	// process input (async)
	input := make(chan string)
	go func(ch chan<- string) {
		for {
			input, err := readInput()
			if err != nil {
				log.Println("error reading input:", err)
				ch <- "ESC"
			}
			ch <- input
		}
	}(input)

	// game loop
	for {
		// process movement
		select {
		case inp := <-input:
			if inp == "ESC" {
				lives = 0
			}
			movePlayer(inp)
		default:
		}
		moveGhosts()

		// process collisions
		for _, g := range ghosts {
			if player.row == g.position.row && player.col == g.position.col {
				ghostsStatusMx.RLock()
				if g.status == GhostStatusNormal {
					lives = lives - 1
					if lives != 0 {
						moveCursor(player.row, player.col)
						fmt.Print(cfg.Death)
						moveCursor(len(maze)+2, 0)
						ghostsStatusMx.RUnlock()
						updateGhosts(ghosts, GhostStatusNormal)
						time.Sleep(1000 * time.Millisecond) //dramatic pause before reseting player position
						player.row, player.col = player.startRow, player.startCol
					}
				} else if g.status == GhostStatusBlue {
					ghostsStatusMx.RUnlock()
					updateGhosts([]*ghost{g}, GhostStatusNormal)
					g.position.row, g.position.col = g.position.startRow, g.position.startCol
				}
			}
		}

		// check game over
		if numDots == 0 || lives == 0 {
			moveCursor(player.row, player.col)
			fmt.Print(cfg.Death)
			moveCursor(len(maze)+2, 0)
			break
		}

		// update screen
		printScreen()

		// repeat
		time.Sleep(200 * time.Millisecond)
	}
}
