package messagehandling

import (
	"fmt"
)

func HandleMessage(username, msg string) {
	fmt.Printf("%v: %v\n", username, msg)
}
