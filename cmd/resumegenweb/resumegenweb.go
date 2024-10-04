package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"resumegenerator/internal/templates"
	"strings"
)

const STATIC_DIR string = "web/static/"

func main() {
	fmt.Println("Hello, Go!")

	t := templates.New()

	http.HandleFunc("GET /dashboard", func(w http.ResponseWriter, r *http.Request) {
		err := t.Render(w, "dashboard", nil)
		if err != nil {
			// send error
			fmt.Println(err.Error())
			return
		}
	})

	http.HandleFunc("GET /{filename}", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		filename := r.PathValue("filename")
		mimeType := "text/plain"

		filenameParts := strings.Split(filename, ".")
		if len(filenameParts) < 1 {
			// send error
			return
		}
		ext := filenameParts[len(filenameParts)-1]

		switch ext {
		case "js":
			mimeType = "application/javascript"
			break
		case "css":
			mimeType = "text/css"
			break
		}

		b, err := os.ReadFile(STATIC_DIR + filename)
		if err != nil {
			// send error
			return
		}

		w.Header().Set("content-type", mimeType)
		w.Write(b)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
