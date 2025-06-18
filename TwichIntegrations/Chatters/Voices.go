package chatters

import (
	"database/sql"
	"math/rand"
)

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

func GetVoice(conn *sql.DB, UID int) Voice {
	// From cache
	var val, exists = VoicesCache[UID]
	if exists {
		return val
	}
	// From DB
	var storedVoice, stored = GetVoiceFromDB(UID, conn)
	if stored {
		VoicesCache[UID] = *storedVoice
		return *storedVoice
	}
	// New
	var newVoice = GenerateRandomVoice()
	newVoice.UserID = UID
	VoicesCache[UID] = newVoice
	RegisterVoice(UID, newVoice, conn) 
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

func GetVoiceFromDB(UID int, conn *sql.DB) (*Voice, bool) {
	var rows, err = conn.Query("select * from Voices where UserID = $1", UID)
	if err != nil {
		return nil, false
	}
	defer rows.Close()
	for rows.Next() {
		var outp = Voice{}
		rows.Scan(&outp.UserID, &outp.Accent, &outp.Speed, &outp.Pitch, &outp.Capital)
		return &outp, true
	}
	return nil, false
}

func RegisterVoice(UID int, toSave Voice, conn *sql.DB) {
	var _, err = conn.Exec("insert into Voices (UserID,	Accent,	Speed, Pitch, Capital) values ($1, $2, $3, $4, $5)", 
		UID, toSave.Accent, toSave.Speed, toSave.Pitch, toSave.Capital)
	if err != nil {
		return
	}
}
