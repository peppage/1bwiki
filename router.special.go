package main

import (
	"net/http"
	"strconv"
	"strings"

	mdl "1bwiki/model"

	"1bwiki/view"

	"github.com/kataras/iris"
)

func edit(c *iris.Context) {
	pageTitle := c.URLParam("title")

	urlTitle, yes := needsRedirect(pageTitle)
	if yes {
		c.Redirect("/special/edit?title="+urlTitle, http.StatusMovedPermanently)
		return
	}

	pv := mdl.GetView(mdl.NameSpace[mdl.WikiPage], pageTitle)

	if pv.NiceTitle == "" {
		pv.NameSpace = ""
		pv.Title = pageTitle
		pv.NiceTitle = strings.Replace(pageTitle, "_", " ", -1)
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
	pageTitle := c.URLParam("title")

	urlTitle, yes := needsRedirect(pageTitle)
	if yes {
		c.Redirect("/special/history?title="+urlTitle, http.StatusMovedPermanently)
		return
	}

	p, _ := strconv.Atoi(c.URLParam("page"))
	revs, err := mdl.GetPageRevisions(pageTitle, p, 50)
	if err != nil {
		c.Error("", http.StatusInternalServerError)
	}
	val := c.Session().Get("user")
	niceTitle := mdl.NiceTitle(pageTitle)
	totalPages := int(mdl.GetAmountOfRevisionsForPage(pageTitle) / 50)
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
	p, err := mdl.GetViews()
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
	t := mdl.GetRandomViewTitle()
	c.Redirect("/pages/"+t, http.StatusTemporaryRedirect)
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
