package actions

import (
	"colossus/models"
	"colossus/pkg"
	"fmt"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/defaults"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/pkg/errors"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/google/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	gu, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	var userInfo = map[string]string{
		"name":       defaults.String(gu.Name, gu.NickName),
		"providerID": gu.UserID,
		"provider":   gu.Provider,
		"email":      gu.Email,
	}
	tx := c.Value("tx").(*pop.Connection)
	u := []models.User{}
	err = tx.Eager("Group").Where("email = ? and provider = ? and provider_id = ?",
		userInfo["email"], userInfo["provider"], userInfo["providerID"]).All(&u)
	if err != nil {
		return errors.WithStack(err)
	}
	var groupName string
	if len(u) > 0 {
		groupName = u[0].Group.Name
	} else {
		u := &models.User{}
		u.Name = userInfo["name"]
		u.Provider = userInfo["provider"]
		u.ProviderID = userInfo["providerID"]
		u.Email = userInfo["email"]
		// TODO: Use cache replace db query group table
		group := models.Group{}
		err := tx.Where("Name = ?", "defaultgroup").First(&group)
		if err != nil {
			return c.Render(http.StatusBadRequest, r.String(err.Error()))
		}
		u.Group = &group
		if err = tx.Save(u); err != nil {
			return errors.WithStack(err)
		}
		groupName = group.Name
	}

	tokenString, err := pkg.CreateJWTToken(userInfo["email"], userInfo["provider"], userInfo["providerID"], groupName)
	if err != nil {
		return errors.WithStack(err)
	}
	var data = map[string]string{
		"jwttoken": tokenString,
	}
	return c.Render(http.StatusOK, r.JSON(data))
}
