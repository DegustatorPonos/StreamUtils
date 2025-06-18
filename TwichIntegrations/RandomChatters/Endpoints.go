package randomchatters

import (
	"fmt"
	"io"
	"net/http"
	"os"

	ev "StreamTTS/EnvVariables"
	twichcomm "StreamTTS/TwichComm"

	"golang.org/x/net/websocket"
)

const ViewsLocation = "RandomChatters/View"

var _WSConnections = []*websocket.Conn{}

func RegisterEndpoints() {
		http.Handle("/api/rnd/ws", websocket.Handler(handleWS))
		http.HandleFunc("/rnd", indexView) 
		http.HandleFunc("/rnd/control", controlView) 
		http.HandleFunc("/api/rnd/connect", connectAPIRequest)
		http.HandleFunc("/api/rnd/disconnect", disconnectAPIRequest)
		http.HandleFunc("/api/rnd/dumpMessage", GetMostRecentMessage)
}

// Checks if the request includes the valid token
func authorizeRequest(r *http.Request) bool {
	var provided = r.URL.Query().Get("token")
	return provided == ev.Enviroment.AppAPIKey
}

func indexView(w http.ResponseWriter, r *http.Request) {
	var ViewPath = fmt.Sprintf("%v/index.html", ViewsLocation)
	var file, fopenerr = os.ReadFile(ViewPath)
	if fopenerr != nil {
		fmt.Fprintf(w, "<h1>An error occured while reading the requested file</h1><p>Original error: %v</p>", fopenerr.Error())
		return
	}
	fmt.Fprint(w, string(file))
}

func controlView(w http.ResponseWriter, r *http.Request) {
	var ViewPath = fmt.Sprintf("%v/control.html", ViewsLocation)
	var file, fopenerr = os.ReadFile(ViewPath)
	if fopenerr != nil {
		fmt.Fprintf(w, "<h1>An error occured while reading the requested file</h1><p>Original error: %v</p>", fopenerr.Error())
		return
	}
	fmt.Fprint(w, string(file))
}
func connectAPIRequest(w http.ResponseWriter, r *http.Request) {
	if !authorizeRequest(r) {
		w.WriteHeader(403)
		return
	}
	var channelInfo, err = twichcomm.GetChannelInfo("physickdev")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	CurrentState.CurrentCahtter = &channelInfo.Data[0]
	onConnect()
}

func disconnectAPIRequest(w http.ResponseWriter, r *http.Request) {
	if !authorizeRequest(r) {
		w.WriteHeader(403)
		return
	}
	CurrentState.CurrentCahtter = nil
	onDisconnect()
}

func GetMostRecentMessage(w http.ResponseWriter, r *http.Request) {
	if(len(CurrentState.Messages) == 0) {
	fmt.Fprint(w, "\n0")
		return
	}
	fmt.Fprintf(w, "%v\n", <-CurrentState.Messages)
	fmt.Fprintf(w, "%d\n", len(CurrentState.Messages))
}

func handleWS(ws *websocket.Conn) {
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
