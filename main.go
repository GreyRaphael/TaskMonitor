package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const (
	username = "admin"
	password = "123456"
)

func main() {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", basicAuth(statusHandler(tmpl)))
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

func statusHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := getNginxState()
		if err := tmpl.Execute(w, state); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getNginxState() nginxState {
	isRunning := nginxStatus()
	if isRunning {
		resp, err := http.Get("http://127.0.0.1/nginx_status")
		if err != nil {
			return nginxState{false, "0", "0", "0"}
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nginxState{false, "0", "0", "0"}
		}

		// Parse the response body
		re := regexp.MustCompile(`\d+`)
		matches := re.FindAllString(string(body), 3)

		return nginxState{
			true,
			matches[0],
			matches[1],
			matches[2],
		}
	}
	return nginxState{false, "0", "0", "0"}
}

func controlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		action := r.FormValue("action")
		switch action {
		case "start":
			startNginx()
		case "stop":
			stopNginx()
		}

		time.Sleep(time.Second) // sleep 1s, then Redirect to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
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

type nginxState struct {
	Running     bool
	ActiveNum   string
	AcceptedNum string
	HandledNum  string
}
