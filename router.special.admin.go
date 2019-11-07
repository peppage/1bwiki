package main

import (
	"net/http"

	mdl "1bwiki/model"
	"1bwiki/view"

	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

func admin(c *iris.Context) {
	val := c.Session().Get("user")
	u := val.(*mdl.User)
	p := &view.AdminPage{
		User: u,
		URL:  "/special/admin",
	}
	view.WritePageTemplate(c.GetRequestCtx(), p)
	c.HTML(http.StatusOK, "")
}

func adminHandle(c *iris.Context) {
	err := mdl.SetAnonEditing(c.FormValueString("anon") == "on")
	if err != nil {
		log.WithError(err).Error("Failed to set anon editing")
		c.Error("Failed to set anon editing", http.StatusInternalServerError)
		return
	}
	err = mdl.SetSignups(c.FormValueString("signup") == "on")
	if err != nil {
		log.WithError(err).Error("Failed to set signups")
		c.Error("Failed to change signups", http.StatusInternalServerError)
	}
	c.Redirect("/special/admin", http.StatusSeeOther)
}
