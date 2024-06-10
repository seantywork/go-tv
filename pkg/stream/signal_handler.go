package stream

import (
	"log"
	"net/http"
	"time"
)

func createSignalHandlerForWS() {

	http.HandleFunc("/signal", signalHandler)

	log.Fatal(http.ListenAndServe(*ADDR, nil))

}

func signalHandler(w http.ResponseWriter, r *http.Request) {

	UPGRADER.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := UPGRADER.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgraded")
		return
	}

	c.SetReadDeadline(time.Time{})

	var sinfo SIGNAL_INFO

	err = c.ReadJSON(&sinfo)

	if err != nil {

		log.Printf("failed to read uinfo: %s\n", err.Error())

		return

	}

	uid := sinfo.UserID

	old_c, okay := USER_SIGNAL[uid]

	if okay {

		log.Printf("uid: %s, exists, removing previous conn\n", uid)

		old_c.Close()

		delete(USER_SIGNAL, uid)

	}

	for k, v := range USER_SIGNAL {

		new_uinfo := SIGNAL_INFO{
			Command: "ADDUSER",
			UserID:  uid,
		}

		v.WriteJSON(&new_uinfo)

		old_uinfo := SIGNAL_INFO{

			Command: "ADDUSER",
			UserID:  k,
		}

		c.WriteJSON(&old_uinfo)

	}

	USER_SIGNAL[uid] = c

	for {

		geninfo := SIGNAL_INFO{}

		c.ReadJSON(&geninfo)

	}

}