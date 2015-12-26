package main

import (
	"net/http"
	"strings"

	mdl "1bwiki/model"
	"1bwiki/tmpl/page"
	"1bwiki/tmpl/special"

	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
)

func edit(c *echo.Context) error {
	n, t := seperateNamespaceAndTitle(c.Query("title"))

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

func history(c *echo.Context) error {
	n, t := seperateNamespaceAndTitle(c.Query("title"))
	urlTitle := convertTitleToUrl(t)
	if urlTitle != t {
		if n != "" {
			n += ":"
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/special/history?title="+n+urlTitle)
	}
	revs, err := mdl.GetPageRevisions(c.Query("title"))
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	niceTitle := mdl.NiceTitle(c.Query("title"))
	return c.HTML(http.StatusOK, page.History(val.(*mdl.User), niceTitle, revs))
}

func recentChanges(c *echo.Context) error {
	revs, err := mdl.GetRevisions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, special.Recentchanges(val.(*mdl.User), revs))
}

func pages(c *echo.Context) error {
	p, err := mdl.GetPages()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, special.Pages(val.(*mdl.User), p))
}

func random(c *echo.Context) error {
	t := mdl.GetRandomPageViewTitle()
	return c.Redirect(http.StatusTemporaryRedirect, "/"+t)
}
