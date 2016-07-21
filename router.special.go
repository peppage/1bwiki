package main

import (
	"net/http"
	"strconv"
	"strings"

	mdl "1bwiki/model"

	"1bwiki/view"

	log "github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
)

func edit(c *iris.Context) {
	n, t := seperateNamespaceAndTitle(c.URLParam("title"))

	ut := strings.ToLower(t)
	if strings.HasPrefix(ut, noEditArea) {
		c.NotFound()
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
		c.Redirect("/special/edit?title="+n+urlTitle, http.StatusMovedPermanently)
		return
	}

	pv := mdl.GetPageView(n, t)

	if pv.NiceTitle == "" {
		if n != "" {
			n += ":"
		}
		pv.NameSpace = n
		pv.Title = t
		pv.NiceTitle = strings.Replace(n+t, "_", " ", -1)
	}
	val := c.Session().Get("user")
	p := &view.ArticleEdit{
		User: val.(*mdl.User),
		Page: pv,
	}
	view.WritePageTemplate(c.GetRequestCtx(), p)
	c.HTML(http.StatusOK, "")
}

func history(c *iris.Context) {
	n, t := seperateNamespaceAndTitle(c.URLParam("title"))
	urlTitle := convertTitleToUrl(t)
	if urlTitle != t {
		if n != "" {
			n += ":"
		}
		c.Redirect("/special/history?title="+n+urlTitle, http.StatusTemporaryRedirect)
		return
	}
	p, _ := strconv.Atoi(c.URLParam("page"))
	revs, err := mdl.GetPageRevisions(c.URLParam("title"), p, 50)
	if err != nil {
		c.Error("", http.StatusInternalServerError)
	}
	val := c.Session().Get("user")
	niceTitle := mdl.NiceTitle(c.URLParam("title"))
	totalPages := int(mdl.GetAmountOfRevisionsForPage(c.URLParam("title")) / 50)
	page := &view.ArticleHistory{
		User:       val.(*mdl.User),
		NiceTitle:  niceTitle,
		Revs:       revs,
		Page:       p,
		TotalPages: totalPages,
	}
	view.WritePageTemplate(c.GetRequestCtx(), page)
	c.HTML(http.StatusOK, "")
}

func recentChanges(c *iris.Context) {
	limit := 50
	if c.URLParam("limit") != "" {
		var err error
		limit, err = c.URLParamInt("limit")
		if err != nil {
			limit = 50
		}
	}
	revs, err := mdl.GetRevisions(limit)
	if err != nil {
		c.EmitError(http.StatusInternalServerError)
		return
	}
	val := c.Session().Get("user")
	p := &view.RecentChangesPage{
		URL:   "special/recentchanges",
		User:  val.(*mdl.User),
		Revs:  revs,
		Limit: limit,
	}
	view.WritePageTemplate(c.GetRequestCtx(), p)
	c.HTML(http.StatusOK, "")
}

func pages(c *iris.Context) {
	p, err := mdl.GetPageViews()
	if err != nil {
		c.EmitError(http.StatusInternalServerError)
	}
	val := c.Session().Get("user")
	page := &view.PageListPage{
		URL:   "special/pages",
		User:  val.(*mdl.User),
		Pages: p,
	}
	view.WritePageTemplate(c.GetRequestCtx(), page)
	c.HTML(http.StatusOK, "")
}

func random(c *iris.Context) {
	t := mdl.GetRandomPageViewTitle()
	c.Redirect("/"+t, http.StatusTemporaryRedirect)
}

func delete(c *iris.Context) {
	val := c.Session().Get("user")
	p := &view.DeletePage{
		PageTitle: c.URLParam("title"),
		URL:       "/special/delete",
		User:      val.(*mdl.User),
	}
	view.WritePageTemplate(c.GetRequestCtx(), p)
	c.HTML(http.StatusOK, "")
}

func deleteHandle(c *iris.Context) {
	val := c.Session().Get("user")
	err := mdl.DeletePage(val.(*mdl.User), c.FormValueString("title"))
	if err != nil {
		c.Error("Failed deleting page", http.StatusInternalServerError)
		return
	}
	c.Redirect("/", http.StatusSeeOther)
}

func users(c *iris.Context) {
	u, err := mdl.GetUsers()
	if err != nil {
		c.Error("Failed to get users", http.StatusInternalServerError)
		return
	}

	val := c.Session().Get("user")
	p := &view.UsersListPage{
		Users: u,
		URL:   "/special/users",
		User:  val.(*mdl.User),
	}
	view.WritePageTemplate(c.GetRequestCtx(), p)
	c.HTML(http.StatusOK, "")
}
