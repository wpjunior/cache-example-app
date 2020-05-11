package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	units "github.com/docker/go-units"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
<html>
<body>
   <a href="/cache/1h">One hour cache</a><br/>
   <a href="/cache/1h/1mb">One hour cache and 10MB de payload</a>
</body>
<html>
`)
	})
	mux.HandleFunc("/cache/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		duration, _ := time.ParseDuration(parts[2])

		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("X-Accel-Expires", strconv.Itoa(int(duration/time.Second)))

		fmt.Fprintf(w, "Cached at %s", time.Now().UTC().String())

		if len(parts) > 3 {
			size, err := units.FromHumanSize(parts[3])
			if err != nil {
				fmt.Fprintf(w, "Err: %s", err.Error())
				return
			}
			fmt.Fprintln(w, "Payload")
			chunks := int(size / 16)
			for i := 0; i < chunks; i++ {
				w.Write([]byte("0123456789ABCDEF"))
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8888", mux))
}
