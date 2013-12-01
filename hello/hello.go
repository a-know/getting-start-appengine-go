package hello

import (
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Greeting struct {
	Author  string
	Content string
	Date    time.Time
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/sign", sign)
}

func root(res http.ResponseWriter, req *http.Request) {
	context := appengine.NewContext(req)
	// ancestor query を発行
	q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(context)).Order("-Date").Limit(10)
	// サイズ10のsliceを作成
	greetings := make([]Greeting, 0, 10)
	if _, err := q.GetAll(context, &greetings); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := guestbookTemplate.Execute(res, greetings); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// 全 Greeting エンティティに設定する親keyを返す
func guestbookKey(context appengine.Context) *datastore.Key {
	// 親keyの作成。stringIDは "default_guestbook" で固定
	return datastore.NewKey(context, "Guestbook", "default_guestbook", 0, nil)
}

func sign(res http.ResponseWriter, req *http.Request) {
	context := appengine.NewContext(req)
	g := Greeting{
		Content: req.FormValue("content"),
		Date:    time.Now(),
	}
	if u := user.Current(context); u != nil {
		g.Author = u.String()
	}
	key := datastore.NewIncompleteKey(context, "Greeting", guestbookKey(context))
	_, err := datastore.Put(context, key, &g)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/", http.StatusFound)
}

var guestbookTemplate = template.Must(template.New("book").Parse(guestbookTemplateHTML))

const guestbookTemplateHTML = `
<html>
  <body>
    {{range .}}
      {{with .Author}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
      <pre>{{.Content}}</pre>
    {{end}}
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`
