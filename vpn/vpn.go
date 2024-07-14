package vpn

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

// https://tunnelblick.net/cAppleScriptSupport.html#controlling-tunnelblick-from-the-command-line

// no need to check the connection state
// if you try to connect to vpn when already connected, nothing happens
// if you try to disconnect when vpn is not connected, nothing happens

// this is how you could do it, if there is need for it
// vpnInstruction := "osascript -e 'tell application \"Tunnelblick\" to get state of first configuration'"
// two possible states: EXITING (i.e., not connected) and CONNECTED

// for work tasks I need VPN connection
func StartOpenVPNConnection(wg *sync.WaitGroup) {
	defer wg.Done()

	vpnInstruction := "osascript -e 'tell application \"Tunnelblick\" to connect \"mpetkovic\"'"

	cmd := exec.Command("bash", "-c", vpnInstruction)

	// err := cmd.Run()
	err := cmd.Start() // don't wait for command to finish
	if err != nil {
		log.Fatalf("Failed to start OpenVPN: %v", err)
	} else {
		fmt.Println("OpenVPN Connection is started")
	}
}

// for study projects I don't want VPN connection
func CloseOpenVPNConnection(wg *sync.WaitGroup) {
	defer wg.Done()

	vpnInstruction := "osascript -e 'tell application \"Tunnelblick\" to disconnect all'"

	cmd := exec.Command("bash", "-c", vpnInstruction)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to close OpenVPN: %v", err)
	} else {
		fmt.Println("OpenVPN Connection is disconnected")
	}
}
