package tor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Start launches the Tor process using the local 'torrc' configuration file.
var cmd *exec.Cmd

// Start launches the Tor process using the local 'torrc' configuration file.
func Start() error {
	torPath := filepath.Join("tor", "tor")
	cmd = exec.Command(torPath, "-f", "torrc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

// WaitForHostname waits for the Tor hidden service hostname file to appear.
// It polls for up to 30 seconds and returns the .onion address when available.
func WaitForHostname() (string, error) {
	hostnamePath := filepath.Join("hidden_service", "hostname")

	for i := 0; i < 30; i++ {
		data, err := os.ReadFile(hostnamePath)
		if err == nil {
			return string(data), nil
		}
		time.Sleep(1 * time.Second)
	}

	return "", fmt.Errorf("timed out waiting for Tor hidden service hostname")
}

// Stop terminates the Tor process if it is running.
func Stop() error {
	if cmd != nil && cmd.Process != nil {
		return cmd.Process.Kill()
	}
	return nil
}

func isChrooted() bool {
	link, err := os.Readlink("/proc/1/root")
	return err == nil && link != "/"
}
