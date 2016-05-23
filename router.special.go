package main

import (
	"net/http"
	"strconv"
	"strings"

	mdl "1bwiki/model"

	"1bwiki/view"

	"github.com/labstack/echo"
	"github.com/peppage/echo-middleware/session"
)

func edit(c echo.Context) error {
	n, t := seperateNamespaceAndTitle(c.QueryParam("title"))

	ut := strings.ToLower(t)
	if strings.HasPrefix(ut, noEditArea) {
		return echo.NewHTTPError(http.StatusForbidden, "Editing of special pages disallowed")
	}

	urlTitle := convertTitleToUrl(t)
	if urlTitle != t {
		if n != "" {
			n += ":"
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/special/edit?title="+n+urlTitle)
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
	session := session.Default(c)
	val := session.Get("user")
	p := &view.ArticleEdit{
		User: val.(*mdl.User),
		Page: pv,
	}
	return c.HTML(http.StatusOK, view.PageTemplate(p))
}

func history(c echo.Context) error {
	n, t := seperateNamespaceAndTitle(c.QueryParam("title"))
	urlTitle := convertTitleToUrl(t)
	if urlTitle != t {
		if n != "" {
			n += ":"
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/special/history?title="+n+urlTitle)
	}
	p, _ := strconv.Atoi(c.QueryParam("page"))
	revs, err := mdl.GetPageRevisions(c.QueryParam("title"), p, 50)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	niceTitle := mdl.NiceTitle(c.QueryParam("title"))
	totalPages := int(mdl.GetAmountOfRevisionsForPage(c.QueryParam("title")) / 50)
	page := &view.ArticleHistory{
		User:       val.(*mdl.User),
		NiceTitle:  niceTitle,
		Revs:       revs,
		Page:       p,
		TotalPages: totalPages,
	}
	return c.HTML(http.StatusOK, view.PageTemplate(page))
}

func recentChanges(c echo.Context) error {
	limit := 50
	if c.QueryParam("limit") != "" {
		var err error
		limit, err = strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			limit = 50
		}
	}
	revs, err := mdl.GetRevisions(limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	p := &view.RecentChangesPage{
		URL:   "special/recentchanges",
		User:  val.(*mdl.User),
		Revs:  revs,
		Limit: limit,
	}
	return c.HTML(http.StatusOK, view.PageTemplate(p))
}

func pages(c echo.Context) error {
	p, err := mdl.GetPageViews()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	page := &view.PageListPage{
		URL:   "special/pages",
		User:  val.(*mdl.User),
		Pages: p,
	}
	return c.HTML(http.StatusOK, view.PageTemplate(page))
}

func random(c echo.Context) error {
	t := mdl.GetRandomPageViewTitle()
	return c.Redirect(http.StatusTemporaryRedirect, "/"+t)
}

func delete(c echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	p := &view.DeletePage{
		PageTitle: c.QueryParam("title"),
		URL:       "/special/delete",
		User:      val.(*mdl.User),
	}
	return c.HTML(http.StatusOK, view.PageTemplate(p))
}

func deleteHandle(c echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	err := mdl.DeletePage(val.(*mdl.User), c.FormValue("title"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

func users(c echo.Context) error {
	u, err := mdl.GetUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	p := &view.UsersListPage{
		Users: u,
		URL:   "/special/users",
		User:  val.(*mdl.User),
	}
	return c.HTML(http.StatusOK, view.PageTemplate(p))
}
