package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

type Survey struct {
	ClassTitle    string
	ClassCode     string
	ClassSchedule string
	TeacherName   string
	Results       map[string]interface{}
	TotalVote     string
}

var (
	resultsSlice [7]map[string]int
	surveys      = make(map[string]*Survey)
	mu           sync.Mutex
)

func main() {
	// 最初に表示する設問を追加
	survey := &Survey{
		ClassTitle:    "総合ゼミナール",
		ClassCode:     "E51R",
		ClassSchedule: "金曜1",
		TeacherName:   "日大　太郎",
		TotalVote:     "0",
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
	http.HandleFunc("/admin/create", createHandler)
	http.HandleFunc("/submit-success1", submitSuccess1Handler)
	http.HandleFunc("/submit-success2", submitSuccess2Handler)

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
			questionKey := fmt.Sprintf("q%d", i) // name属性の名前を用意、q1からq7のname属性を用意
			response := r.FormValue(questionKey) // name属性に対応するvalue属性をresponseに代入
			resultsSlice[i-1][response]++        // i-1番目のマップでresponseキーに対応するvalueを加算
		}

		// 投票人数をカウントしてstringに型変換しフィールドを更新
		totalVotes := 0
		for _, results := range resultsSlice {
			for _, count := range results {
				totalVotes += count
			}
		}
		survey.TotalVote = fmt.Sprintf("%d", totalVotes/7)

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

		http.Redirect(w, r, "/submit-success1", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("survey.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, survey)
}

func submitSuccess1Handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "surveySuccess.html")
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

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		survey := &Survey{
			ClassTitle:    r.FormValue("classTitle"),
			ClassCode:     r.FormValue("classCode"),
			ClassSchedule: r.FormValue("classSchedule"),
			TeacherName:   r.FormValue("teacherName"),
			TotalVote:     "0",
		}
		surveys[r.FormValue("classTitle")] = survey
		http.Redirect(w, r, "/submit-success2", http.StatusSeeOther)

		return
	}
	http.ServeFile(w, r, "create.html")
}

func submitSuccess2Handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "createSuccess.html")
}
