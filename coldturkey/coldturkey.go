package coldturkey

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

var instruction = "\"/Applications/Cold Turkey Blocker.app/Contents/MacOS/Cold Turkey Blocker\" -start work_session"

func StartWorkSession(wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("bash", "-c", instruction)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("❌ [cold-turkey] Failed to start Work Session: %v", err)
	} else {
		fmt.Println("✔️ [cold-turkey] started Work Session")
	}
}
