package main

import (
	"net/http"
	"strconv"
	"strings"

	mdl "1bwiki/model"
	"1bwiki/tmpl/page"
	"1bwiki/tmpl/special"

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
	return c.HTML(http.StatusOK, page.Editpage(val.(*mdl.User), pv))
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
	return c.HTML(http.StatusOK, page.History(val.(*mdl.User),
		niceTitle, revs, p, totalPages))
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
	return c.HTML(http.StatusOK, special.Recentchanges(val.(*mdl.User), revs, limit))
}

func pages(c echo.Context) error {
	p, err := mdl.GetPageViews()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, special.Pages(val.(*mdl.User), p))
}

func random(c echo.Context) error {
	t := mdl.GetRandomPageViewTitle()
	return c.Redirect(http.StatusTemporaryRedirect, "/"+t)
}

func delete(c echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, special.Delete(val.(*mdl.User), c.QueryParam("title")))
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
	//session := session.Default(c)
	//val := session.Get("user")
	//s.Users(u)
	p := &view.UsersPage{
		Users: u,
		URL:   "/special/users",
	}
	return c.HTML(http.StatusOK, view.PageTemplate(p))
}
