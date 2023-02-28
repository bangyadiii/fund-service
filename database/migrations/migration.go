package migrations

import (
	"backend-crowdfunding/database"
	"backend-crowdfunding/src/model"
)

type Migration struct {
	DB *database.DB
}

func (m *Migration) RunMigration() error {
	return m.DB.AutoMigrate(
		&model.User{},
		&model.Campaign{},
		&model.CampaignImage{},
		&model.Transaction{},
	)
}

func (m *Migration) DropDatabase() error {
	return m.DB.Migrator().DropTable(
		&model.User{},
		&model.Campaign{},
		&model.CampaignImage{},
		&model.Transaction{},
	)
}
