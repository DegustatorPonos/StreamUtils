package chatters

import (
	"database/sql"
	"fmt"
)

type Chatter struct {
	Id int
	Username string
	ELO int
}

// Username to ID
var ChattersIDCache = make(map[string]int)

const StartELO int = 128

// Returns user ID from BD or -1 if he is not written
func GetChatterID(username string, conn *sql.DB) int {
	var val, exists = ChattersIDCache[username]
	if exists {
		return val
	}
	var rows, err = conn.Query("select * from Chatters where Username like $1", username)
	if err != nil {
		fmt.Println(err.Error())
		return -1
	}
	defer rows.Close()
	for rows.Next() {
		var chatter = Chatter{Id: -1}
		var err = rows.Scan(&chatter.Id, &chatter.Username, &chatter.ELO)
		if err != nil {
			continue
		}
		ChattersIDCache[username] = chatter.Id
		return chatter.Id 
	}
	return -1
}

func RegisterChatter(username string, conn *sql.DB) {
	var _, err = conn.Exec("insert into Chatters (Username, ELO) values ($1, $2)", username, StartELO)
	if err != nil {
		return
	}
}
