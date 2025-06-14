package tor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var cmd *exec.Cmd

func Start() error {
	torPath := filepath.Join("tor", "tor")
	cmd = exec.Command(torPath, "-f", "torrc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func WaitForHostname() (string, error) {
	hostnamePath := filepath.Join("hidden_service", "hostname")
	for i := 0; i < 30; i++ {
		data, err := ioutil.ReadFile(hostnamePath)
		if err == nil {
			return string(data), nil
		}
		time.Sleep(1 * time.Second)
	}
	return "", fmt.Errorf("timed out waiting for Tor hidden service hostname")
}

func Stop() error {
	if cmd != nil && cmd.Process != nil {
		return cmd.Process.Kill()
	}
	return nil
}
