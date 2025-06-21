package messagehandling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

const RateDeclineCoeff = 5

type userRecord struct {
	Username string `json:"username"`
	Value float64 `json:"value"`
	LastMessageTime time.Time `json:"lastmessagetime"`
}

type activityMeter struct {
	Metrics map[string]userRecord `json:"metrics"`
}

var ActivityMeter activityMeter 

func registerMessage(username string, _ string) {
	var val, exists = ActivityMeter.Metrics[username]
	if !exists {
		ActivityMeter.Metrics[username] = userRecord{
			Username: username,
			Value: 1,
			LastMessageTime: time.Now(),
		}
		return
	}
	val.recalculateMetric()
	val.LastMessageTime = time.Now()
	ActivityMeter.Metrics[username] = val
}

func Init() {
	ActivityMeter = activityMeter{
		Metrics: make(map[string]userRecord),
	}
	RegisterHandler(&Handler{
		Condition: func(_, _ string) bool { return true }, 
		Action: registerMessage,
	})
}

func (data *userRecord) recalculateMetric() {
	var dt = time.Now().Unix() - data.LastMessageTime.Unix() + 1
	data.Value = (data.Value / float64(dt)) + 1
}

func RegisterEndpoints() {
	http.HandleFunc("/api/metrics/raw", getMetricsRaw)
	http.HandleFunc("/api/metrics/list", getUsersTierlist)
}

func getUsersTierlist(w http.ResponseWriter, r *http.Request) {
	type temp struct {
		Data []string `json:"data"`
	}
	var outp = temp{Data: GetUserActivityRating()}
	var body, jsonerr = json.Marshal(outp)
	if jsonerr != nil {
		fmt.Printf("An error occured while creating user activity tierlist. \nOriginal message: %v\n", jsonerr.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(body)
}

func getMetricsRaw(w http.ResponseWriter, r *http.Request) {
	var resp = activityMeter{
		Metrics: make(map[string]userRecord),
	}
	for k, v := range ActivityMeter.Metrics {
		var clone = userRecord{
			Username: v.Username,
			Value: v.Value,
			LastMessageTime: v.LastMessageTime,
		}
		clone.recalculateMetric()
		resp.Metrics[k] = clone
	}
	var body, jsonerr = json.Marshal(resp)
	if jsonerr != nil {
		fmt.Printf("An error occured while creating raw metrics resp. \nOriginal message: %v\n", jsonerr.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(body)
}

// Returns the results sorted
func GetUserActivityRating() []string {
	var users = make([]userRecord, 0, len(ActivityMeter.Metrics))
	for _, v := range ActivityMeter.Metrics {
		users = append(users, v)
		users[len(users) - 1].recalculateMetric()
	}
	sort.Slice(users, func(i, j int)bool { return users[i].Value < users[j].Value })
	var outp = make([]string, 0, len(users))
	for _, v := range users {
		outp = append(outp, v.Username)
	}
	return outp
}
