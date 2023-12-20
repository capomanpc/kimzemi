package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"sync"
)

// アンケートのデータを格納するマップ
var surveys = make(map[string][]string)
var surveysLock sync.RWMutex

func main() {
	http.HandleFunc("/create", createSurvey)
	http.HandleFunc("/survey/", takeSurvey)
	http.HandleFunc("/result/", showResult)

	http.ListenAndServe(":8080", nil)
}

func createSurvey(w http.ResponseWriter, r *http.Request) {
	var surveyTitle string

	if r.Method == http.MethodPost {
		// フォームから題名を取得
		surveyTitle = r.FormValue("title")
		if surveyTitle != "" {
			// アンケートデータをマップに格納
			surveysLock.Lock()
			surveys[surveyTitle] = []string{"Option 1", "Option 2", "Option 3"}
			surveysLock.Unlock()
		}
	}

	// アンケート作成フォームの表示
	tmpl, err := template.ParseFiles("templates/create.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// アンケートの題名をテンプレートに渡す
	tmpl.Execute(w, surveyTitle)
}

func takeSurvey(w http.ResponseWriter, r *http.Request) {
	// URLからアンケートの題名を抽出
	surveyTitle := r.URL.Path[len("/survey/"):]

	// アンケートデータをマップから取得
	surveysLock.RLock()
	options, exists := surveys[surveyTitle]
	surveysLock.RUnlock()
	if !exists {
		http.NotFound(w, r)
		return
	}

	// アンケートテンプレートを読み込み
	tmpl, err := template.ParseFiles("templates/survey.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// アンケートを表示
	tmpl.Execute(w, options)
}

func showResult(w http.ResponseWriter, r *http.Request) {
	// URLからアンケートの題名を抽出
	surveyTitle := r.URL.Path[len("/result/"):]

	// アンケートデータをマップから取得
	surveysLock.RLock()
	options, exists := surveys[surveyTitle]
	surveysLock.RUnlock()
	if !exists {
		http.NotFound(w, r)
		return
	}

	// アンケート結果を計算
	results := make([]struct {
		Option string
		Count  int
	}, len(options))

	for i, option := range options {
		// ここで実際のアンケート結果のカウント処理を行う
		// この例ではランダムな数を生成しています
		count := rand.Intn(100)
		results[i] = struct {
			Option string
			Count  int
		}{
			Option: option,
			Count:  count,
		}
	}

	// 集計結果を表示するテンプレートに渡す
	tmpl, err := template.ParseFiles("templates/results.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Results []struct {
			Option string
			Count  int
		}
	}{
		Title:   surveyTitle,
		Results: results,
	}

	// テンプレートを実行して結果を出力
	tmpl.Execute(w, data)
}
