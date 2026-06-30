package helpFunc

import (
	"bufio"
	"fmt"
	"os"
)

func ReplLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		inp := CleanInput(scanner.Text())
		if len(inp) == 0 {
			continue
		}

		if val, ok := commands[inp[0]]; ok {
			if inp[0] == "explore" {
				setEploreName(inp)
			}
			if inp[0] == "catch" {
				setCatchName(inp)
			}
			if inp[0] == "inspect" {
				setInspectName(inp)
			}
			val.callback()
		} else {
			fmt.Println("Unknown command")
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("error encounter in the scanner")
		}
	}

}
