package orb

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

func StartOrbStack(wg *sync.WaitGroup) {
	defer wg.Done()

	var cmd = exec.Command("bash", "-c", "pgrep -x 'OrbStack' > /dev/null")

	err := cmd.Run()
	if err != nil {
		cmd = exec.Command("bash", "-c", "open -j -a OrbStack")
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Failed to start OrbStack: %v", err)
		}
		fmt.Println("✔️ [orbstack] started orbstack")
	} else {
		fmt.Println("✔️ [orbstack] orbstack is already running")
	}
}
