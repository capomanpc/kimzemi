package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

type Survey struct {
	ClassTitle  string
	ClassCode   string
	ClassTiming string
	TeacherName string
	Results     map[string]interface{}
}

var (
	resultsSlice [7]map[string]int
	surveys      = make(map[string]*Survey)
	mu           sync.Mutex
)

func main() {
	// 最初に表示する設問を追加
	survey := &Survey{
		ClassTitle:  "総合ゼミナール",
		ClassCode:   "E51R",
		ClassTiming: "金曜1",
		TeacherName: "日大　太郎",
	}
	surveys["総合ゼミナール"] = survey

	for i := range resultsSlice {
		resultsSlice[i] = make(map[string]int) // mapを初期化
	}

	for i := 1; i <= 7; i++ {
		for j := 1; j <= 5; j++ {
			key := fmt.Sprintf("op%d", j)
			resultsSlice[i-1][key] = 0 // スライスに格納されているマップのvalueを0に初期化
		}
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/survey/", surveyHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/admin/results/", resultsHandler)
	http.HandleFunc("/admin/results/getdata/", getDataHandler)
	//http.HandleFunc("/admin/create", createHandler)

	fmt.Println("サーバー起動 : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, surveys)
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

		for i := 1; i <= 7; i++ {
			questionKey := fmt.Sprintf("q%d", i) // 質問キーを生成 (例: "q1", "q2", ...)
			response := r.FormValue(questionKey) // フォームからの回答を取得
			resultsSlice[i-1][response]++
		}

		var resultsQ1 = resultsSlice[0]
		var resultsQ2 = resultsSlice[1]
		var resultsQ3 = resultsSlice[2]
		var resultsQ4 = resultsSlice[3]
		var resultsQ5 = resultsSlice[4]
		var resultsQ6 = resultsSlice[5]
		var resultsQ7 = resultsSlice[6]

		survey.Results = map[string]interface{}{
			"resultsQ1": resultsQ1,
			"resultsQ2": resultsQ2,
			"resultsQ3": resultsQ3,
			"resultsQ4": resultsQ4,
			"resultsQ5": resultsQ5,
			"resultsQ6": resultsQ6,
			"resultsQ7": resultsQ7,
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("survey.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, survey)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "admin" && password == "password" {
			http.Redirect(w, r, "/admin", http.StatusFound)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	} else {
		// リクエストメソッドがGETの場合
		// ログイン画面を表示
		http.ServeFile(w, r, "login.html")
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, surveys)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/admin/results/"):] // URLからtitleを読み取る
	mu.Lock()
	survey, ok := surveys[title] //titleに合ったアンケートデータが入った構造体をsurveyに代入
	mu.Unlock()

	if !ok {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("results.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, survey)
}

func getDataHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/admin/results/getdata/"):] // URLからtitleを読み取る
	mu.Lock()
	survey, ok := surveys[title] //titleに合ったアンケートデータが入った構造体をsurveyに代入
	mu.Unlock()

	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(survey.Results)
}
