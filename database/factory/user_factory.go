package factory

import (
	"backend-crowdfunding/src/model"
	"time"

	"github.com/bxcodec/faker/v3"
)

func UserFactory() *model.User {
	return &model.User{
		ID:              IDGenerator.Generate(),
		Name:            faker.Name(),
		Occupation:      faker.Word(),
		Email:           faker.Email(),
		Password:        "$2y$10$lZT8tZ0krPVrcwYh1fm4eOuxKJ.DFOu4XhrFvjiR.gB5UaK9SoY2y", // password
		Avatar:          faker.URL(),
		IsGoogleAccount: false,
		Role:            model.RoleAdmin,
		Transactions:    nil,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
