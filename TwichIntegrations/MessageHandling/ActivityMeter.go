package messagehandling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	ev "StreamTTS/EnvVariables"
)

const RateDeclineCoeff float64 = float64(1) / 50

type userRecord struct {
	Username string `json:"username"`
	Weight float64 `json:"value"`
	LastMessageTime time.Time `json:"lastmessagetime"`
	TotalMessages uint64 `json:"totalmessages"`
}

type ActivityMeter struct {
	Metrics map[string]userRecord `json:"metrics"`
	WeigthsSum float64 `json:"weigthssum"`
}

var ActivityMeterState ActivityMeter 

func registerMessage(username string, _ string) {
	var val, exists = ActivityMeterState.Metrics[username]
	if !exists {
		ActivityMeterState.Metrics[username] = userRecord{
			Username: username,
			Weight: 1,
			LastMessageTime: time.Now(),
			TotalMessages: 1,
		}
		ActivityMeterState.WeigthsSum++
		return
	}
	val.Weight++
	val.TotalMessages++
	val.LastMessageTime = time.Now()
	var delta = val.recalculateMetric()
	ActivityMeterState.Metrics[username] = val
	ActivityMeterState.WeigthsSum += delta + 1
}

func Init() {
	ActivityMeterState = ActivityMeter{
		Metrics: make(map[string]userRecord),
		WeigthsSum: 0,
	}
	RegisterHandler(&Handler{
		Condition: func(_, _ string) bool { return true }, 
		Action: registerMessage,
	})
}

func (data *userRecord) recalculateMetric() float64 {
	var dt = time.Now().Unix() - data.LastMessageTime.Unix() + 1
	var delta = -1 * data.Weight
	data.Weight -= float64(dt) * RateDeclineCoeff
	if data.Weight < 0 {
		data.Weight = 0
	}
	delta += data.Weight
	return delta
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
	var probe = ProbeUserActivity()
	var body, jsonerr = json.Marshal(probe)
	if jsonerr != nil {
		fmt.Printf("An error occured while creating raw metrics resp. \nOriginal message: %v\n", jsonerr.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(body)
}

// Returns the results sorted
func GetUserActivityRating() []string {
	var users = make([]userRecord, 0, len(ActivityMeterState.Metrics))
	for _, v := range ActivityMeterState.Metrics {
		users = append(users, v)
		users[len(users) - 1].recalculateMetric()
	}
	sort.Slice(users, func(i, j int)bool { return users[i].Weight < users[j].Weight })
	var outp = make([]string, 0, len(users))
	for _, v := range users {
		outp = append(outp, v.Username)
	}
	return outp
}

func ProbeUserActivity() *ActivityMeter {
	if !ev.Config.ActivityMetrics {
		return &ActivityMeter{}
	}
	var outp = ActivityMeter{
		Metrics: make(map[string]userRecord),
		WeigthsSum: ActivityMeterState.WeigthsSum,
	}
	for k, v := range ActivityMeterState.Metrics {
		var clone = userRecord{
			Username: v.Username,
			Weight: v.Weight,
			LastMessageTime: v.LastMessageTime,
			TotalMessages: v.TotalMessages,
		}
		var delta = clone.recalculateMetric()
		outp.Metrics[k] = clone
		outp.WeigthsSum += delta
	}
	return &outp
}
