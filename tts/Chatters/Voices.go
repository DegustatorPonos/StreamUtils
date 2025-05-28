package chatters

import "math/rand"

type Voice struct {
	UserID  int
	Accent  string
	Speed   int
	Pitch   int
	Capital int
}

var VoicesCache = make(map[int]Voice, 0)

var Voices = []string{
	"an",
	"bg",
	"bn",
	"bs",
	"ca",
	"cs",
	"cy",
	"da",
	"de",
	"el",
	"en",
	"en-gb",
	"en-sc",
	"en-wm",
	"eo",
	// TODO: finish
}

func GetVoice(UID int) Voice {
	var val, exists = VoicesCache[UID]
	if exists {
		return val
	}
	var newVoice = GenerateRandomVoice()
	newVoice.UserID = UID
	VoicesCache[UID] = newVoice
	return newVoice
}

func GenerateRandomVoice() Voice {
	var accId = rand.Intn(len(Voices))
	var speed = rand.Intn(100) + 12
	var pitch = rand.Intn(100) + 1
	var cap = rand.Intn(200) + 1
	return Voice{
		Accent: Voices[accId],
		Speed: speed,
		Pitch: pitch,
		Capital: cap,
	}
}
