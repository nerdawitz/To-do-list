package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"text/template"
)

type Task struct {
	Todo     string
	Reminder string
	Priority string
}

var tasks []Task = []Task{
	Task{
		Todo:     "Here you can see your tasks!",
		Reminder: "Date & Time will appear here",
		Priority: "Low",
	},
	Task{
		Todo:     "Task #2",
		Reminder: "Date & Time will appear here.",
		Priority: "Medium",
	},
}

func main() {

	fs := http.FileServer(http.Dir("."))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Starting server...")

	url := "http://localhost:8080"

	err := openURL(url)
	if err != nil {
		fmt.Println("Error opening URL:", err)
	}

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tpl := template.Must(template.ParseFiles("index.html"))
		tpl.Execute(w, tasks)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		task := r.PostFormValue("task")
		reminder := r.PostFormValue("reminder")
		priority := r.PostFormValue("priority")

		reminderDate := reminder[0:10]
		reminderTime := reminder[11:]
		reminderFormatted := reminderDate + ", " + reminderTime

		tpl := template.Must(template.ParseFiles("index.html"))
		tpl.ExecuteTemplate(w, "task-list-element", Task{Todo: task, Reminder: reminderFormatted, Priority: priority})
	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		tpl := template.Must(template.ParseFiles("index.html"))
		tpl.Execute(w, nil)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-task/", h2)
	http.HandleFunc("/del-task/", h3)
	http.ListenAndServe(":8080", nil)
}

func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {

	case "darwin":
		cmd = exec.Command("open", url)

	case "linux":
		cmd = exec.Command("xdg-open", url)

	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)

	default:
		return fmt.Errorf("Unsupported Platform! Couldn't open browser! open the browser, and navigate to 'localhost:8080'")
	}

	return cmd.Start()
}
