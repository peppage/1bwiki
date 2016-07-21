package main

import (
	"encoding/gob"
	"net/http"
	"strings"
	"time"

	mdl "1bwiki/model"
	"1bwiki/setting"
	"1bwiki/view"

	log "github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
)

const noEditArea = "special"

func convertTitleToUrl(t string) string {
	firstChar := string(t[0])
	t = strings.ToUpper(firstChar) + string(t[1:])
	t = strings.Replace(t, "%20", "_", -1)
	t = strings.Replace(t, " ", "_", -1)
	return t
}

func seperateNamespaceAndTitle(t string) (namespace string, title string) {
	URL := strings.Trim(t, "/")
	if strings.Contains(URL, ":") {
		split := strings.Split(URL, ":")
		namespace = split[0]
		title = split[1]
	} else {
		title = URL
	}
	return namespace, title
}

func root(c *iris.Context) {
	c.Redirect("/pages/Main_Page", http.StatusMovedPermanently)
}

func wikiPage(c *iris.Context) {
	n, t := seperateNamespaceAndTitle(c.Param("name"))
	log.WithFields(log.Fields{
		"namespace": n,
		"title":     t,
	}).Debug("separating namespace and title")

	ul := strings.ToLower(c.Param("name"))
	if strings.HasPrefix(ul, "/"+noEditArea) {
		c.EmitError(http.StatusForbidden)
		return
	}

	urlTitle := convertTitleToUrl(t)

	log.WithFields(log.Fields{
		"urlTitle": urlTitle,
		"t":        t,
	}).Debug("Should this page be redirected?")
	if urlTitle != t {
		if n != "" {
			n += ":"
		}
		c.Redirect("/pages/"+n+urlTitle, http.StatusMovedPermanently)
		return
	}

	if c.URLParam("oldid") != "" {
		pv, err := mdl.GetPageVeiwByID(c.URLParam("oldid"))
		if err != nil {
			c.EmitError(http.StatusInternalServerError)
			return
		}

		val := c.Session().Get("user")
		if c.URLParam("diff") != "" {
			pv2, err := mdl.GetPageVeiwByID(c.URLParam("diff"))
			if err != nil {
				c.EmitError(http.StatusInternalServerError)
				return
			}
			p := &view.ArticleDiff{
				User:  val.(*mdl.User),
				Page:  pv,
				Page2: pv2,
			}
			view.WritePageTemplate(c.GetRequestCtx(), p)
			c.HTML(http.StatusOK, "")
			return
		}
		p := &view.ArticleOld{
			User: val.(*mdl.User),
			Page: pv,
		}
		view.WritePageTemplate(c.GetRequestCtx(), p)
		c.HTML(http.StatusOK, "")
		return
	}

	pv := mdl.GetPageView(n, t)

	if pv.NiceTitle != "" && !pv.Deleted {
		val := c.Session().Get("user")
		p := &view.Article{
			User: val.(*mdl.User),
			Page: pv,
		}
		view.WritePageTemplate(c.GetRequestCtx(), p)
		c.HTML(http.StatusOK, "")
		return
	}
	if n != "" {
		n += ":"
	}
	c.Redirect("/special/edit?title="+n+t, http.StatusTemporaryRedirect)
}

func savePage(c *iris.Context) {
	val := c.Session().Get("user")
	u, ok := val.(*mdl.User)
	if !ok {
		log.WithFields(log.Fields{
			"user": u,
		}).Error("User saving page is invalid")
		c.EmitError(http.StatusBadRequest)
		return
	}

	minor := c.FormValueString("minor") == "on"
	p, err := mdl.CreateOrUpdatePage(u, mdl.CreatePageOptions{
		Title:     c.FormValueString("title"),
		Namespace: c.FormValueString("namespace"),
		Text:      c.FormValueString("text"),
		Comment:   c.FormValueString("summary"),
		IsMinor:   minor,
	})
	if err != nil {
		c.EmitError(http.StatusBadRequest)
		return
	}
	c.Redirect("/pages/"+p.Title, http.StatusSeeOther)
}

func init() {
	gob.Register(&mdl.User{})
	setting.Initialize()
	ll, err := log.ParseLevel(setting.LogLevel)
	if err == nil {
		log.SetLevel(ll)
	}

	iris.Config.Sessions.Cookie = "id"
	iris.Config.Sessions.Expires = time.Hour * 48
	iris.Config.Sessions.GcDuration = time.Duration(2) * time.Hour
	iris.Config.Gzip = true

}

func main() {
	mdl.SetupDb()

	iris.Use(&sessionMiddleware{})
	iris.Static("/static", "./static", 1)
	iris.Get("/", root)
	iris.Get("/pages/*name", wikiPage)

	special := iris.Party("/special")
	special.Get("/edit", edit)
	special.Post("/edit", savePage)
	special.Get("/history", history)
	special.Get("/recentchanges", recentChanges)
	special.Get("/pages", pages)
	special.Get("/users", users)
	special.Get("/register", register)
	special.Post("/register", registerHandle)
	special.Get("/login", login)
	special.Post("/login", loginHandle)
	special.Get("/logout", logout)
	special.Get("/random", random)
	special.Get("/delete", delete)
	special.Post("/delete", deleteHandle)

	user := iris.Party("/preferences")
	user.Use(&loggedInMiddleware{})
	user.Get("", prefs)
	user.Get("/password", prefsPasword)
	user.Post("/password", handlePrefsPassword)

	a := iris.Party("/admin")
	a.Use(&adminMiddleware{})
	a.Get("", admin)
	a.Post("", adminHandle)

	iris.Listen(":" + setting.HttpPort)
}
