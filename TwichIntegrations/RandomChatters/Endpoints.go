package randomchatters

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

const ViewsLocation = "RandomChatters/View"

var _WSConnections = []*websocket.Conn{}

func RegisterEndpoints() {
		http.HandleFunc("/rnd", Index) 
		http.Handle("/api/rnd/ws", websocket.Handler(AcceptWS))
		http.HandleFunc("/api/rnd/dumpMessage", GetMostRecentMessage)
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

func AcceptWS(ws *websocket.Conn) {
	_WSConnections = append(_WSConnections, ws)
	var buf = make([]byte, 1024)
	for {
		var _, err = ws.Read(buf)
		if err == io.EOF {
			break
		}
		ws.Write(buf)
	}
	ws.Close()
}
