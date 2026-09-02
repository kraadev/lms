package main

import (
	"os"
	"os/exec"
)

// Main entrypoint redirects to cmd/api
func main() {
	cmd := exec.Command("go", "run", "./cmd/api")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}
