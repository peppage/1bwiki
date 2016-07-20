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
	"github.com/labstack/echo"
	"github.com/peppage/echo-middleware/session"
)

var store session.CookieStore

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

	if urlTitle != t {
		if n != "" {
			n += ":"
		}
		c.Redirect("/"+n+urlTitle, http.StatusMovedPermanently)
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

func savePage(c echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	u, ok := val.(*mdl.User)
	if !ok {
		log.WithFields(log.Fields{
			"user": u,
		}).Error("User saving page is invalid")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid User")
	}

	minor := c.FormValue("minor") == "on"
	p, err := mdl.CreateOrUpdatePage(u, mdl.CreatePageOptions{
		Title:     c.FormValue("title"),
		Namespace: c.FormValue("namespace"),
		Text:      c.FormValue("text"),
		Comment:   c.FormValue("summary"),
		IsMinor:   minor,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Save page  failed")
	}
	return c.Redirect(http.StatusSeeOther, "/"+p.Title)
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
}

func main() {
	mdl.SetupDb()

	iris.Use(&sessionMiddleware{})
	iris.Static("/static", "./static", 1)
	iris.Get("/", root)
	iris.Get("/pages/*name", wikiPage)

	iris.Listen(":" + setting.HttpPort)
}
