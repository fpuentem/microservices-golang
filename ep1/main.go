package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello World")
		d, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ops", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Hello %s", d)
	})

	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("Goodbye")
	})

	http.ListenAndServe(":9090", nil)
}

// curl -v -d 'Kelly' localhost:9090
// curl -v localhost:9090/goodbey
// https://www.youtube.com/watch?v=VzBGi_n65iU&list=PLmD8u-IFdreyh6EUfevBcbiuCKzFk0EW_
