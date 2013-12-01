package hello

import (
	"fmt"
	"net/http"

	"appengine"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(res http.ResponseWriter, req *http.Request) {
	// returns an appengine.Context value associated with the current request
	context := appengine.NewContext(req)
	u := user.Current(context)
	if u == nil {
		url, err := user.LoginURL(context, req.URL.String())
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusFound)
		return
	}
	fmt.Fprintf(res, "こんにちは, %vさん！", u)
}
