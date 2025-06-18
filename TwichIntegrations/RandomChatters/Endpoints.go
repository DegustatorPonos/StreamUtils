package randomchatters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	ev "StreamTTS/EnvVariables"
	twichcomm "StreamTTS/TwichComm"

	"golang.org/x/net/websocket"
)

const ViewsLocation = "RandomChatters/View"

var _WSConnections = []*websocket.Conn{}

func RegisterEndpoints() {
		http.Handle("/api/rnd/ws", websocket.Handler(handleWS))

		http.HandleFunc("/rnd/view", indexView) 
		http.HandleFunc("/rnd/control", controlView) 

		http.HandleFunc("/api/rnd/connect", connectAPIRequest)
		http.HandleFunc("/api/rnd/disconnect", disconnectAPIRequest)
		http.HandleFunc("/api/rnd/dumpMessage", GetMostRecentMessage)
		http.HandleFunc("/api/rnd/bannedusers", ignoredChatterAPIRequest) 
		http.HandleFunc("/api/rnd/ban", banAPIRequest) 
		http.HandleFunc("/api/rnd/pardon", pardonAPIRequest) 

		http.HandleFunc("/rnd/style.css", cssEndpoint) 
		http.HandleFunc("/rnd/control.js", controlScriptEndpoint) 
		http.HandleFunc("/rnd/view.js", viewScriptEndpoint) 
}

// Checks if the request includes the valid token
func authorizeRequest(r *http.Request) bool {
	var provided = r.URL.Query().Get("token")
	return provided == ev.Enviroment.AppAPIKey
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

func ignoredChatterAPIRequest(w http.ResponseWriter, r *http.Request) {
	if !authorizeRequest(r) {
		w.WriteHeader(403)
		return
	}
	type list struct {
		Chatters []string `json:"chatters"`
	}
	var outp = list{Chatters: IgnoredChatters}
	var body, err = json.Marshal(outp)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(body)
}

func banAPIRequest(w http.ResponseWriter, r *http.Request) {
	if !authorizeRequest(r) {
		w.WriteHeader(403)
		return
	}
	var user = r.URL.Query().Get("user")
	if user == "" || slices.Contains(IgnoredChatters, user) {
		return
	}
	IgnoredChatters = append(IgnoredChatters, user)
}

func pardonAPIRequest(w http.ResponseWriter, r *http.Request) {
	if !authorizeRequest(r) {
		w.WriteHeader(403)
		return
	}
	var user = r.URL.Query().Get("user")
	if user == "" {
		return
	}
	var index = slices.Index(IgnoredChatters, user)
	if index != -1 {
		IgnoredChatters = append(IgnoredChatters[:index], IgnoredChatters[index+1:]...)
	}
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
	if CurrentState != nil && CurrentState.CurrentCahtter != nil {
		var event = ConnectEvent{
			Type: "connect",
			UserName: CurrentState.CurrentCahtter.DisplayName,
			UserPfp: CurrentState.CurrentCahtter.ProfileImageUrl,
		}
		var payload, _ = json.Marshal(event)
		ws.Write(payload)
	}
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

func serveFile(w http.ResponseWriter, r *http.Request, fileName string) {
	var filePath = fmt.Sprintf("%v/%v", ViewsLocation, fileName)
	var file, err = os.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	fmt.Fprint(w, string(file))
}

func indexView(w http.ResponseWriter, r *http.Request) {
	serveFile(w, r, "index.html")
}

func controlScriptEndpoint(w http.ResponseWriter, r *http.Request) {
	serveFile(w, r, "control.js")
}

func viewScriptEndpoint(w http.ResponseWriter, r *http.Request) {
	serveFile(w, r, "view.js")
}

func cssEndpoint(w http.ResponseWriter, r *http.Request) {
	serveFile(w, r, "style.css")
}

func controlView(w http.ResponseWriter, r *http.Request) {
	serveFile(w, r, "control.html")
}
