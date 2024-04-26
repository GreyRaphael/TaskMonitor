package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	http.HandleFunc("/", statusHandler)
	http.HandleFunc("/control", controlHandler)
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	status := nginxStatus()
	fmt.Println("Nginx status:", status)
	if status {
		fmt.Fprintln(w, "Nginx is running")
	} else {
		fmt.Fprintln(w, "Nginx is stopped")
	}
}

func controlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "start" {
			startNginx()
			fmt.Fprintln(w, "Nginx started")
		} else if action == "stop" {
			stopNginx()
			fmt.Fprintln(w, "Nginx stopped")
		}
	} else {
		http.Error(w, "Invalid request method.", 405)
	}
}

func nginxStatus() bool {
	processName := "Notepad2.exe"
	cmd := exec.Command("tasklist", "/fi", fmt.Sprintf("imagename eq %s", processName))

	if output, err := cmd.Output(); err == nil {
		return strings.Contains(string(output), processName)
	}
	return false
}

func startNginx() {
	exec.Command("cmd", "/c", "start", "nginx").Run()
}

func stopNginx() {
	exec.Command("nginx", "-s", "stop").Run()
}
