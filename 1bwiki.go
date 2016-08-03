package main

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	mdl "1bwiki/model"
	"1bwiki/setting"
	"1bwiki/view"

	log "github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
)

// NeedsRedirect checks that the page title is properly formatted
// returns if the page needs a redirect and the proper page name
func needsRedirect(title string) (string, bool) {
	firstChar := string(title[0])
	t := strings.ToUpper(firstChar) + string(title[1:])
	t = strings.Replace(t, "%20", "_", -1)
	t = strings.Replace(t, " ", "_", -1)
	return t, title != t
}

func root(c *iris.Context) {
	c.Redirect("/pages/Main_Page", http.StatusMovedPermanently)
}

func showDiffPage(c *iris.Context) {
	oldPage, err := mdl.GetViewByID(c.URLParam("oldid"))
	if err != nil {
		c.EmitError(http.StatusInternalServerError)
		return
	}

	val := c.Session().Get("user")
	diffPage, err := mdl.GetViewByID(c.URLParam("diff"))
	if err != nil {
		c.EmitError(http.StatusInternalServerError)
		return
	}
	p := &view.ArticleDiff{
		User:  val.(*mdl.User),
		Page:  oldPage,
		Page2: diffPage,
	}
	view.WritePageTemplate(c.GetRequestCtx(), p)
	c.HTML(http.StatusOK, "")
}

func showOldPage(c *iris.Context) {
	pv, err := mdl.GetViewByID(c.URLParam("oldid"))
	if err != nil {
		c.EmitError(http.StatusInternalServerError)
		return
	}
	val := c.Session().Get("user")
	p := &view.ArticleOld{
		User: val.(*mdl.User),
		Page: pv,
	}
	view.WritePageTemplate(c.GetRequestCtx(), p)
	c.HTML(http.StatusOK, "")
}

func wikiPage(c *iris.Context) {
	pageTitle := strings.Trim(c.Param("name"), "/")

	urlTitle, yes := needsRedirect(pageTitle)
	if yes {
		c.Redirect("/pages/"+urlTitle, http.StatusMovedPermanently)
		return
	}

	if c.URLParam("oldid") != "" && c.URLParam("diff") != "" {
		showDiffPage(c)
		return
	}

	if c.URLParam("oldid") != "" {
		showOldPage(c)
		return
	}

	// Showing regular page
	pv := mdl.GetView("", pageTitle)

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

	// Page doesn't exist redirect to edit
	c.Redirect("/special/edit?title="+pageTitle, http.StatusTemporaryRedirect)
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
	err := mdl.CreateOrUpdatePage(u, mdl.CreatePageOptions{
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
	c.Redirect("/pages/"+c.FormValueString("title"), http.StatusSeeOther)
}

func init() {
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
	iris.Get("/static/*filename", func(c *iris.Context) {
		filename := "static" + c.Param("filename")
		data, err := Asset(filename)
		if err != nil {
			log.Error("asset problem" + err.Error())
			c.NotFound()
			return
		}
		fi, _ := AssetInfo(filename)
		_ = c.ServeContent(bytes.NewReader(data), c.Param("filename"), fi.ModTime(), true)
	})

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

	preferences := special.Party("/preferences")
	preferences.Use(&loggedInMiddleware{})
	preferences.Get("", prefs)
	preferences.Post("", handlePrefs)
	preferences.Get("/appearence", prefsAppearence)
	preferences.Post("/appearence", handlePrefsAppearence)
	preferences.Get("/password", prefsPasword)
	preferences.Post("/password", handlePrefsPassword)

	administration := special.Party("/admin")
	administration.Use(&adminMiddleware{})
	administration.Get("", admin)
	administration.Post("", adminHandle)

	iris.Listen(":" + setting.HttpPort)
}
