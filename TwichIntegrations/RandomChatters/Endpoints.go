package randomchatters

import (
	"fmt"
	"net/http"
	"os"
)

const ViewsLocation = "RandomChatters/View"

func RegisterEndpoints() {
		http.HandleFunc("/rnd", Index) 
		http.HandleFunc("/api/rnd/getCurrentMessage", GetMostRecentMessage) 
}

func Index(w http.ResponseWriter, r *http.Request) {
	var ViewPath = fmt.Sprintf("%v/index.html", ViewsLocation)
	var file, fopenerr = os.ReadFile(ViewPath)
	if fopenerr != nil {
		fmt.Fprintf(w, "<h1>An error occured while reading the requested file</h1><p>Original error: %v</p>", fopenerr.Error())
		return
	}
	fmt.Fprint(w, string(file))
}

func GetMostRecentMessage(w http.ResponseWriter, r *http.Request) {
	if(len(CurrentState.Messages) == 0) {
	fmt.Fprint(w, "\n0")
		return
	}
	fmt.Fprintf(w, "%v\n", <-CurrentState.Messages)
	fmt.Fprintf(w, "%d\n", len(CurrentState.Messages))
}
