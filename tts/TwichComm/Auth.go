package twichcomm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	ev "StreamTTS/EnvVariables"
)

const AuthLink string = 
`https://id.twitch.tv/oauth2/authorize
?response_type=code
&client_id=%v
&redirect_uri=http://localhost:3000/auth
&scope=user:read:chat`

const OAuthLink = "https://id.twitch.tv/oauth2/token"

const ValidationLink = "https://id.twitch.tv/oauth2/validate"

type OAuthResp struct {
	Access_token string `json:"access_token"`
	Expires_in int `json:"expires_in"`
	Refresh_token string `json:"refresh_token"`
	Scope []string  `json:"scope"`
	Token_type string  `json:"token_type"`
}

// User authentication will come to this channel and it will be closed after
var AuthKeyChan chan string 

// Checks a stored value of 
func AuthenticateApp() bool { 
	AuthKeyChan = make(chan string, 1) 
	var oauthResp, err = ExchangeCode()
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		printAuthRequest(ev.Enviroment.TwichAPIKey)
		return false
	}
	ev.Enviroment.UserToken = oauthResp.Access_token
	if IsTokenValid() {
		AuthKeyChan <- ev.Enviroment.UserCode
		return true
	}
	printAuthRequest(ev.Enviroment.TwichAPIKey)
	return false
}

func printAuthRequest(ApiKey string) {
	fmt.Println("Go to this link to get an auth token:")
	var AuthURI = fmt.Sprintf(AuthLink, ApiKey)
	fmt.Println(strings.ReplaceAll(AuthURI, "\n", ""));
}

func AuthKeyHttpEndpoint(w http.ResponseWriter, r *http.Request) {
	var code = r.FormValue("code")
	fmt.Printf("User code: %v\n", code)
	ev.Enviroment.UserCode = code
	var token, err = ExchangeCode()
	if err != nil {
		return
	}
	AuthKeyChan <- token.Access_token
	fmt.Fprint(w, "<script>close()</script>") // so that the tab closes by itself
}

// Warning - may take a long time
func IsTokenValid() bool {
	// fmt.Printf("Validating a token \"%v\"\n", ev.Enviroment.UserToken)
	var client = &http.Client{}
	var req, reqerr = http.NewRequest("GET", ValidationLink, nil)
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %v", ev.Enviroment.UserToken))
	if reqerr != nil {
		return false
	}
	var resp, senderr = client.Do(req)
	if senderr != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

// Calls for twich OAuth services
func ExchangeCode() (*OAuthResp, error) {
	var client = &http.Client{}
	var body = getExcangeBody()
	var req, reqerr = http.NewRequest("POST", OAuthLink, bytes.NewReader(body))
	if reqerr != nil {
		return nil, reqerr
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	var outp = OAuthResp{}
	var respBody = parseResponce(resp)
	fmt.Print(string(respBody))
	var jsonErr = json.Unmarshal(respBody, &outp)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &outp, nil
}

func parseResponce(r *http.Response) []byte {
	var buf, err = io.ReadAll(r.Body)
	if err != nil {
		return []byte{}
	}
	return buf
}

func getExcangeBody() []byte {
	var args = map[string]string {
		"client_id": ev.Enviroment.TwichAPIKey,
		"client_secret": ev.Enviroment.TwichAPISecret,
		"code": ev.Enviroment.UserCode,
		"grant_type": "authorization_code",
		"redirect_uri": "http://localhost:3000/auth",
	}
	var params []string = make([]string, 0, 5)
	for k, v := range args {
		params = append(params, fmt.Sprint(k, "=", v))
	}
	var outp = strings.Join(params, "&")
	return []byte(outp)
}

func AddAuthHeaders(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %v", ev.Enviroment.UserToken))
	req.Header.Set("Client-Id", ev.Enviroment.TwichAPIKey)
	req.Header.Set("Content-Type", "application/json")
}
