package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Survey 構造体を定義
type Survey struct {
	Title     string
	Questions []string
	Results   map[string]int
}

var surveys = make(map[string]Survey)

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		title := r.FormValue("title")
		// アンケートの題名をキーとして Survey 構造体を作成
		survey := Survey{
			Title:     title,
			Questions: []string{"Question 1", "Question 2", "Question 3"}, // 質問を追加できます
			Results:   make(map[string]int),
		}
		surveys[title] = survey
		http.Redirect(w, r, fmt.Sprintf("/survey/%s", title), http.StatusSeeOther)
		return
	}

	// create.htmlのテンプレートを表示
	tmpl, _ := template.New("create").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Create Survey</title>
	</head>
	<body>
		<h1>Create a New Survey</h1>
		<form method="post">
			<label for="title">Survey Title:</label>
			<input type="text" id="title" name="title" required>
			<input type="submit" value="Create">
		</form>
	</body>
	</html>
	`)

	tmpl.Execute(w, nil)
}

func surveyHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/survey/"):]

	// アンケートの題名を元にアンケート入力ページを表示
	survey := surveys[title]
	tmpl, _ := template.New("survey").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>{{.Title}} Survey</title>
	</head>
	<body>
		<h1>{{.Title}} Survey</h1>
		<form method="post" action="/submit/{{.Title}}">
			{{range .Questions}}
				<label for="{{.}}">{{.}}:</label>
				<input type="text" id="{{.}}" name="{{.}}" required><br>
			{{end}}
			<input type="submit" value="Submit">
		</form>
	</body>
	</html>
	`)

	tmpl.Execute(w, survey)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/submit/"):]

	if r.Method == http.MethodPost {
		r.ParseForm()
		// アンケートの題名に対応する Survey 構造体から回答を取得
		survey := surveys[title]
		for _, question := range survey.Questions {
			answer := r.FormValue(question)
			// 回答を結果に追加
			survey.Results[answer]++
		}
		surveys[title] = survey
	}

	// アンケート結果出力ページにリダイレクト
	http.Redirect(w, r, fmt.Sprintf("/results/%s", title), http.StatusSeeOther)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/results/"):]

	// アンケート結果を表示
	survey := surveys[title]

	tmpl, _ := template.New("results").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>{{.Title}} Results</title>
	</head>
	<body>
		<h1>{{.Title}} Results</h1>
		<h2>Answers:</h2>
		<table>
			<tr>
				<th>Answer</th>
				<th>Count</th>
			</tr>
			{{range $answer, $count := .Results}}
				<tr>
					<td>{{$answer}}</td>
					<td>{{$count}}</td>
				</tr>
			{{end}}
		</table>
	</body>
	</html>
	`)

	tmpl.Execute(w, survey)
}

func main() {
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/survey/", surveyHandler)
	http.HandleFunc("/submit/", submitHandler)
	http.HandleFunc("/results/", resultsHandler)

	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
