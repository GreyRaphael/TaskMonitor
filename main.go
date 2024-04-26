package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

var (
	username = "admin"
	password = "123456"
)

func main() {
	http.HandleFunc("/", basicAuth(statusHandler))
	http.HandleFunc("/control", basicAuth(controlHandler))
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	var html string
	if nginxStatus() {
		html = `<html><body><p>Nginx is running.</p><form action="/control" method="post"><input type="submit" name="action" value="stop"/></form></body></html>`
	} else {
		html = `<html><body><p>Nginx is stopped.</p><form action="/control" method="post"><input type="submit" name="action" value="start"/></form></body></html>`
	}
	fmt.Fprintln(w, html)
}

func controlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "start" {
			startNginx()
		} else if action == "stop" {
			stopNginx()
		}
		http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to the home page
	} else {
		http.Error(w, "Invalid request method.", 405)
	}
}

func nginxStatus() bool {
	processName := "nginx.exe"
	cmd := exec.Command("tasklist", "/fi", fmt.Sprintf("imagename eq %s", processName))

	if output, err := cmd.Output(); err == nil {
		return strings.Contains(string(output), processName)
	}
	return false
}

func startNginx() {
	exec.Command("cmd", "/c", "start", "./nginx.exe").Run()
}

func stopNginx() {
	exec.Command("taskkill", "/F", "/IM", "nginx.exe").Run()
}
