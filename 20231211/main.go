package main

import (
	"html/template"
	"net/http"
	"sync"
)

type Survey struct {
	Title     string
	Questions []string
	Results   map[string]int
}

var (
	mu      sync.Mutex
	surveys = make(map[string]*Survey)
)

func main() {
	// HTTPハンドラーを設定
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/survey/", surveyHandler)
	http.HandleFunc("/results/", resultsHandler)

	// ポート8080でHTTPサーバーを起動
	http.ListenAndServe(":8080", nil)
}

// メソッドがPOSTだったときにはフォームの入力をSurvey構造体に入れて、その構造体のアドレスをsurveysマップに格納
// surveysマップのkeyは入力されたタイトルでvalueは構造体のアドレス
// メソッドがPOST以外のときは

func createHandler(w http.ResponseWriter, r *http.Request) {
	// httpパッケージで定義されているMethodPostという定数が"POST"
	//　定数を使用しているため打ち間違えを防止できる
	if r.Method == http.MethodPost {
		// フォームからタイトルを取得
		//http.Request構造体のメソッド
		//POSTメソッドに含まれるフォームデータを解析してhttp.Requestのformフィールドにmapとして格納
		r.ParseForm()
		// r.FormValue("title")は、リクエストオブジェクトのメソッドで、Formフィールドから指定したキー("title")に対応する値を取得する
		title := r.FormValue("title")
		// ダミーの質問を定義
		questions := []string{"Question 1", "Question 2", "Question 3"}

		mu.Lock()
		// 新しいSurvey構造体を作成し、サーバー全体で管理されるsurveysマップに追加
		//それぞれのフィールドにtitle、quetions変数を指定、Resultsフィールドに空のmapを作成
		survey := &Survey{
			Title:     title,
			Questions: questions,
			Results:   make(map[string]int),
		}
		surveys[title] = survey //構造体のアドレス変数surveyをsurveysマップのvalueに格納
		mu.Unlock()

		// 作成したアンケートページにリダイレクト
		// クライアントにリダイレクトレスポンスを送信する関数
		// リダイレクトレスポンスとは、クライアントに別のURLに移動するように指示するレスポンス

		// w:レスポンスライターと呼ばれるオブジェクトで、レスポンスを書き込むために使われる
		// レスポンスライターは、io.Writerインターフェースを実装した任意の型であれば使用することができる

		// r: リクエストオブジェクトで、クライアントからのリクエストの情報を保持している
		// リクエストオブジェクトは、net/httpパッケージで定義されているRequest構造体のポインタである

		// 第四引数: ここにはリダイレクトに関するステータスコードが入る
		// ステータスコードを指定することでリダイレクトに関する操作を指定できる
		// 例えばhttp.StatusseeOtherは303というリダイレクトに関するコードである
		// goでは303を直接記述するのではなくhttp.StatusseeOtherという定数を指定する
		//　303についての解説
		/*
			303HTTPステータスコードはPOSTリクエスト後、ブラウザに新しいページを表示させる際に特に使われます。
			通常ユーザーがウェブフォームにデータを送信した後、サーバーはその情報を処理し結果を伝える新しいページにリダイレクトします。
			この際に303ステータスコードが使用され、ブラウザは自動的に新しいURLにGETリクエストを行いページ内容を取得します。
			303ステータスコードの主な目的は、クライアントに対して直接的なレスポンスを返すのではなく、別のURLにアクセスするよう指示することです。
			データ送信後に再読み込みを避ける、データの二重送信を防ぐ際に役立ちます。
			例えば、オンラインショッピングサイトでの商品購入やアンケートの回答、
			問い合わせフォームの送信後に303ステータスコードを利用して、「送信完了」や「購入完了」などのページにユーザーを誘導することが考えられます。
		*/
		//　300番台についての解説
		/*
			301: Moved Permanently - リクエストされたリソースが恒久的に別のURLに移動したことを示します。
			302: Found - リクエストされたリソースが一時的に別のURLに移動したことを示します。
			303: See Other - リクエストされたリソースを別のURLでGETリクエストを使用して取得するようクライアントに指示することを示します。
			307: Temporary Redirect - リクエストされたリソースが一時的に別のURLに移動したことを示しますが、リクエストのメソッドは変更しないことを示します。
			308: Permanent Redirect - リクエストされたリソースが恒久的に別のURLに移動したことを示しますが、リクエストのメソッドは変更しないことを示します。
		*/

		http.Redirect(w, r, "/survey/"+title, http.StatusSeeOther)
		return // ハンドラ関数を終了
	}

	http.ServeFile(w, r, "templates/create.html")

	/*
		このコードは適切ではない為コメントアウト
		// アンケート作成ページを表示するテンプレートを読み込み、レスポンスに書き込む
		// ParseFilesメソッドでテンプレートファイル(create.html)を構文解析してテンプレート構造体を生成
		// ExcuteメソッドでResponseWriterと何らかのデータを引数として解析済みテンプレートにデータを組み込みHTMLをResponseWriterに送る
		// その後ResponseWriterがHTTPレスポンスを書く
		tmpl, err := template.ParseFiles("templates/create.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// wがResponseWriterでnilは...
		tmpl.Execute(w, nil)
	*/
}

// surveyHandlerはアンケートページを処理するためのハンドラーです。
func surveyHandler(w http.ResponseWriter, r *http.Request) {
	// URLからアンケートのタイトルを取得
	/*
		リクエストメッセージに含まれているURLのパスを

		r.URL.Pathについて
		*http.repuest型構造体のオブジェクトrのURLフィールドのPathフィールド
		Pathフィールドとはurlのホスト名の後に続く文字列のこと
		以下の例だと/audio/list.phpの部分であり、?以降のクエリ文字列は含まれない
		http://www.example.com/audio/list.php?page=2&sort_by=price

		以下のように構造体が定義されている

		type Request struct {
		    ...
		    URL *URL
		    ...
		}

		type URL struct {
		    ...
		    Path   string
		    ...
		}

		r.URL.Path[len("/survey/"):]について

		st := "sample"
		fmt.Println(st[3:]) //ple

		上の例のように
		String型変数[3:]
		の構文で文字列の3文字目より後をスライスに格納することができる

	*/
	title := r.URL.Path[len("/survey/"):]
	mu.Lock()
	// タイトルを使って該当のアンケートを取得
	// surveysマップのtitleというkeyに対応するvalueをsurveyに代入
	// valueが存在しない場合にokにfalse代入される
	survey, ok := surveys[title]
	mu.Unlock()

	if !ok {
		// アンケートが存在しない場合は404エラーを返す
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		// POSTリクエストの場合、フォームから受け取った回答を処理\
		// リクエストメッセージのフォームデータを解析し、formフィールドにmapとして格納
		r.ParseForm()
		// survey.Questionsはスライス、_は要素番号でquestionは要素のこと
		// 設問の数だけ繰り返す
		for _, question := range survey.Questions {
			// FormValueでname属性に対応するvalue属性を代入
			answer := r.FormValue(question)
			mu.Lock()
			// 回答をカウントし、surveyのResultsマップに格納
			//Resultsマップはそのまま表示される
			survey.Results[question+" - "+answer]++
			mu.Unlock()
		}

		// 回答が完了したら結果ページにリダイレクト
		http.Redirect(w, r, "/results/"+title, http.StatusSeeOther)
		return
	}

	// アンケートページを表示するテンプレートを読み込み、レスポンスに書き込む
	tmpl, err := template.ParseFiles("templates/survey.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, survey)
}

// resultsHandlerはアンケート結果ページを処理するためのハンドラーです。
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	// URLからアンケートのタイトルを取得
	title := r.URL.Path[len("/results/"):]
	mu.Lock()
	// タイトルを使って該当のアンケートを取得
	survey, ok := surveys[title]
	mu.Unlock()

	if !ok {
		// アンケートが存在しない場合は404エラーを返す
		http.NotFound(w, r)
		return
	}

	// アンケート結果ページを表示するテンプレートを読み込み、レスポンスに書き込む
	tmpl, err := template.ParseFiles("templates/results.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, survey)
}
