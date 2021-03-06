package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	// NOTE: tmpfiles does NOT exist on my server
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", fs))

	http.HandleFunc("/", ServeTemplate)

	fmt.Println("Listening...")
	err := http.ListenAndServe(GetPort(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}

// Get the Port from the environment so we can run on Heroku
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func ServeTemplate(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("templates", "layout.html") // templates/layout.html
	fp := path.Join("templates", r.URL.Path)    // templates/r.URL.Path[1:]

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	// STEP 1: create a new template - looks like it's automatically created
	// STEP 2: parse the string into the template
	// STEP 3: execute the template

	templates, _ := template.ParseFiles(lp, fp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	templates.ExecuteTemplate(w, "layout", nil)
}
