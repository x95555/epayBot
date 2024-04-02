package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func task() {
	startTime := time.Now()

	for {
		time.Sleep(time.Second * 10)
		startTime = startTime.Add(-10 * time.Second)

		url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("User"), os.Getenv("Password"), os.Getenv("Host"), os.Getenv("Port"), os.Getenv("Database"))
		db, err := sql.Open("mysql", url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var m Msg
		m.Key = os.Getenv("Key")

		o1 := payOrder(startTime, db)
		o2 := getBill(startTime, db)
		m.Data = append(m.Data, o1...)
		m.Data = append(m.Data, o2...)

		startTime = time.Now()

		go sendOrder(m)
		_ = db.Close()

	}
}

func sendOrder(m Msg) {

	j, err := json.Marshal(m)
	if err != nil {
		return
	}

	_, _ = http.Post(os.Getenv("PushUrl"), "application/json", bytes.NewReader(j))

}
