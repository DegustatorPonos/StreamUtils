package twichcomm

import (
	"fmt"
	"net/http"
	"strings"
)

const AuthLink string = 
`https://id.twitch.tv/oauth2/authorize
?response_type=code
&client_id=%v
&redirect_uri=http://localhost:3000/auth
&scope=user:read:chat`

// User authentication will come to this channel and it will be closed after
var AuthKeyChan chan string 

func PrintAuthRequest(ApiKey string) {
	fmt.Println("Go to this link to get an auth token:")
	var AuthURI = fmt.Sprintf(AuthLink, ApiKey)
	fmt.Println(strings.ReplaceAll(AuthURI, "\n", ""));
	AuthKeyChan = make(chan string, 1) 
}

func AuthKeyHttpEndpoint(w http.ResponseWriter, r *http.Request) {
	var code = r.FormValue("code")
	AuthKeyChan <- code
}
