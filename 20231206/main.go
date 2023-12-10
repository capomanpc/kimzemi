package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"
)

var (
	results      = make(map[string]int)
	mutex        = &sync.Mutex{}
	toCopyKeysQ1 = []string{}
	resultsQ1    = make(map[string]int)
	toCopyKeysQ2 = []string{}
	resultsQ2    = make(map[string]int)
	toCopyKeysQ3 = []string{}
	resultsQ3    = make(map[string]int)
	toCopyKeysQ4 = []string{}
	resultsQ4    = make(map[string]int)
	toCopyKeysQ5 = []string{}
	resultsQ5    = make(map[string]int)
	toCopyKeysQ6 = []string{}
	resultsQ6    = make(map[string]int)
	toCopyKeysQ7 = []string{}
	resultsQ7    = make(map[string]int)
)

func main() {
	http.HandleFunc("/", serveForm)
	http.HandleFunc("/submit", submitForm)
	http.HandleFunc("/results", showResults)
	http.HandleFunc("/data", getData)

	//mapを0で初期化
	for i := 1; i <= 7; i++ {
		for j := 1; j <= 5; j++ {
			key := fmt.Sprintf("q%dop%d", i, j)
			results[key] = 0
		}
	}

	fmt.Println("サーバー起動 : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	//name属性に対応するvalue属性をstring型変数に代入
	q1 := r.FormValue("q1")
	q2 := r.FormValue("q2")
	q3 := r.FormValue("q3")
	q4 := r.FormValue("q4")
	q5 := r.FormValue("q5")
	q6 := r.FormValue("q6")
	q7 := r.FormValue("q7")

	//resultsマップの対応するkeyのvalueを加算
	mutex.Lock()
	results["q1"+q1]++
	results["q2"+q2]++
	results["q3"+q3]++
	results["q4"+q4]++
	results["q5"+q5]++
	results["q6"+q6]++
	results["q7"+q7]++
	mutex.Unlock()
	http.Redirect(w, r, "/results", http.StatusFound) // "/results"にリダイレクト
}

func showResults(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// resultsマップを分けるためのスライスの生成
func setKeySlice(toCopyKeys *[]string, qNum int) {
	for i := 1; i <= 5; i++ {
		key := fmt.Sprintf("q%dop%d", qNum, i)
		*toCopyKeys = append(*toCopyKeys, key)
	}
}

// 質問ごとにresultsマップを作り直す
func toCopyVal(toCopyKeys []string, newMap map[string]int) {
	for _, key := range toCopyKeys {
		if val, exists := results[key]; exists {
			newMap[key] = val
		}
	}
}

func getData(w http.ResponseWriter, r *http.Request) {

	setKeySlice(&toCopyKeysQ1, 1)
	toCopyVal(toCopyKeysQ1, resultsQ1)

	setKeySlice(&toCopyKeysQ2, 2)
	toCopyVal(toCopyKeysQ2, resultsQ2)

	setKeySlice(&toCopyKeysQ3, 3)
	toCopyVal(toCopyKeysQ3, resultsQ3)

	setKeySlice(&toCopyKeysQ4, 4)
	toCopyVal(toCopyKeysQ4, resultsQ4)

	setKeySlice(&toCopyKeysQ5, 5)
	toCopyVal(toCopyKeysQ5, resultsQ5)

	setKeySlice(&toCopyKeysQ6, 6)
	toCopyVal(toCopyKeysQ6, resultsQ6)

	setKeySlice(&toCopyKeysQ7, 7)
	toCopyVal(toCopyKeysQ7, resultsQ7)

	resultsAll := map[string]interface{}{
		"resultsQ1": resultsQ1,
		"resultsQ2": resultsQ2,
		"resultsQ3": resultsQ3,
		"resultsQ4": resultsQ4,
		"resultsQ5": resultsQ5,
		"resultsQ6": resultsQ6,
		"resultsQ7": resultsQ7,
	}

	// データをJSON形式でクライアントに送信
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultsAll)
}

/*
	template.ParseFiles("template.html")について
	・テンプレートファイルにデータを挿入するために事前に解析(パース)をする関数
	・テンプレートファイルとは静的なHTML構造を持ちつつ、動的な内容を埋め込むための
	　プレースホルダー（変数や特別なタグ）を含んだファイルのこと

	tmpl.Execute(w, results)について
	・解析されたテンプレートをWebページにレンダリングする関数
	・レンダリングとはプレースホルダに実際のデータを入れて最終的なHTMLファイルを生成すること
*/
