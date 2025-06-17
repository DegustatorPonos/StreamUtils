package randomchatters

import (
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var RandomChatter = GetRandomChatterID()
	fmt.Fprintf(w, "<h1>User name: %v </h1>", RandomChatter.UserName)
}
