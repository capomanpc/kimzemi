package main

import (
	"html/template"
	"net/http"
	"sync"
)

// アンケートごとに構造体でまとめて管理
type Survey struct {
	Title     string
	Questions []string
	Results   map[string]int
}

var (
	mu      sync.Mutex
	surveys = make(map[string]*Survey) // アンケートごとにある構造体をまとめるマップ
)

func main() {
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/survey/", surveyHandler)
	http.HandleFunc("/results/", resultsHandler)

	http.ListenAndServe(":8080", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		title := r.FormValue("title")
		questions := []string{"Question 1", "Question 2", "Question 3"}

		mu.Lock()
		survey := &Survey{
			Title:     title,
			Questions: questions,
			Results:   make(map[string]int), //実際のマップを割り当てる
		}
		surveys[title] = survey
		mu.Unlock()

		http.Redirect(w, r, "/survey/"+title, http.StatusSeeOther)
		:return
	}

	http.ServeFile(w, r, "templates/create.html")
}

func surveyHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/survey/"):] // URLからtitleを読み取る
	mu.Lock()
	survey, ok := surveys[title] //titleに合ったアンケートデータが入った構造体をsurveyに代入
	mu.Unlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		// 入力されたアンケート結果を取得し、resultsフィールドに格納
		for _, question := range survey.Questions {
			answer := r.FormValue(question)
			mu.Lock()
			survey.Results[question+" - "+answer]++
			mu.Unlock()
		}

		http.Redirect(w, r, "/results/"+title, http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/survey.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, survey)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/results/"):]
	mu.Lock()
	survey, ok := surveys[title]
	mu.Unlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/results.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, survey)
}
