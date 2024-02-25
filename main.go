package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type FoodEntry struct {
	Name     string
	Calories int
	Protein  int
	Fat      int
	Carbs    int
}

type DailyReport struct {
	Date    string // dd-mm-yyyy
	Entries []FoodEntry
}

var dailyReports = []DailyReport{}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("POST /add", addHandler)
	mux.HandleFunc("/report", reportHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	foodName := r.FormValue("name")
	calories, err := strconv.Atoi(r.FormValue("calories"))
	if err != nil {
		http.Error(w, "Calories must be a number", http.StatusBadRequest)
		return
	}
	protein, err := strconv.Atoi(r.FormValue("protein"))
	if err != nil {
		http.Error(w, "Protein must be a number", http.StatusBadRequest)
		return
	}
	fat, err := strconv.Atoi(r.FormValue("fat"))
	if err != nil {
		http.Error(w, "Fat must be a number", http.StatusBadRequest)
		return
	}
	carbs, err := strconv.Atoi(r.FormValue("carbs"))
	if err != nil {
		http.Error(w, "Carbs must be a number", http.StatusBadRequest)
		return
	}

	foodEntry := FoodEntry{
		Name:     foodName,
		Calories: calories,
		Protein:  protein,
		Fat:      fat,
		Carbs:    carbs,
	}

	// get the current date as dd-mm-yyyy
	currentDate := time.Now().Format("02-01-2006")
	var dailyEntry *DailyReport
	for i, log := range dailyReports {
		if log.Date == currentDate {
			dailyEntry = &dailyReports[i]
			break
		}
	}

	if dailyEntry == nil {
		dailyEntry = &DailyReport{
			Date: currentDate,
			Entries: []FoodEntry{
				foodEntry,
			},
		}

		dailyReports = append(dailyReports, *dailyEntry)
	} else {
		dailyEntry.Entries = append(dailyEntry.Entries, foodEntry)
	}

	http.Redirect(w, r, "/report", http.StatusFound)
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/daily-report.html")

	templateData := struct {
		DailyReports []DailyReport
	}{
		DailyReports: dailyReports,
	}

	fmt.Println(dailyReports)

	t.Execute(w, templateData)
}
