package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/browse/", http.StatusPermanentRedirect)
}

func handleBrowse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	browse(w, r)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	upload(w, r)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	delete(w, r)
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sysInfo(w, r)
}

func main() {
	// Static content.
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/browse/", handleBrowse)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/sysinfo", handleInfo)
	http.HandleFunc("/", handleRoot)

	log.Printf("Listening on port %s...", serverPort)

	// Let's go!
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil))
}
