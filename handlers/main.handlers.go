package handlers

import (
	"net/http"

	"github.com/bmg-c/product-diary/localization"
	"github.com/bmg-c/product-diary/util"
	"github.com/bmg-c/product-diary/views"
)

func NewMainHandler() *MainHandler {
	return &MainHandler{}
}

type MainHandler struct{}

func (mh *MainHandler) HandleLocale(w http.ResponseWriter, r *http.Request) {
	l := util.InitHTMLHandler(w, r)
	var code int = http.StatusOK
	var out []byte
	defer util.RespondHTTP(w, &code, &out)

	locale, _ := util.GetLocaleCookieValue(r)

	util.RenderComponent(&out, views.LocaleSelect(l, locale), r)
}

func (mh *MainHandler) HandleSetLocale(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	util.SetLocaleCookie(w, localization.LocaleFromString(r.Form.Get("locale")))

	w.Header().Set("HX-Redirect", r.Header.Get("Referer"))
}
