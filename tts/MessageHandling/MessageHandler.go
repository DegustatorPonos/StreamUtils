package messagehandling

import (
	"fmt"
	"os/exec"
)

func HandleMessage(username, msg string) {
	fmt.Printf("%v: %v\n", username, msg)
	go SayMsg(msg)
}

func SayMsg(msg string) {
	var cmd = exec.Command("espeak", fmt.Sprintf("\"%v\"", msg), "&")
	cmd.Run()
}
