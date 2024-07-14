package sketchybar

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

func LoadWorkConfiguration(wg *sync.WaitGroup) {
	defer wg.Done()

	var cmd = exec.Command("bash", "-c", "sketchybar --reload sketchybarrc.work")

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to render sketchybar work configuration: %v", err)
	} else {
		fmt.Println("✔️ [sketchybar] switch to work mode")
	}
}
