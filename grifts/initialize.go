package grifts

import (
	"colossus/models"

	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("initialize", func() {
	grift.Desc("initializeGroup", "Create a default group")
	grift.Add("initializeGroup", func(c *grift.Context) error {
		exists, err := models.DB.Where("name = ?", "defaultgroup").Exists("groups")
		if err != nil {
			return err
		}
		if !exists {
			g := &models.Group{Name: "defaultgroup"}
			err := models.DB.Create(g)
			if err != nil {
				return err
			}
		}
		exists, err = models.DB.Where("name = ?", "admin").Exists("groups")
		if err != nil {
			return err
		}
		if !exists {
			g := &models.Group{
				Name: "admin",
			}
			err := models.DB.Create(g)
			return err
		}
		return nil
	})
})
