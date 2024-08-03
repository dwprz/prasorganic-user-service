package util

import (
	"github.com/dwprz/prasorganic-user-service/src/model/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserTest struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserTest(db *gorm.DB, l *logrus.Logger) *UserTest {
	return &UserTest{
		db:     db,
		logger: l,
	}
}

func (u *UserTest) Create() *entity.User {
	query := `
	INSERT INTO 
		users(email, full_name, password, refresh_token) 
	VALUES
		('johndoe@gmail.com', 'John Doe', 'rahasia', 'eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxNzIwMDUsImlkIjoiMV9pUGtNbjk4c19ObXNRZ1Q1TWtlIiwiaXNzIjoicHJhc29yZ2FuaWMtYXV0aC1zZXJ2aWNlIn0.cVJL1ivJ5wDECYwBQtA39R_HMkEaG4HiRHxZSJBl0EL5_EcuKq5v7QscveiFYd7CEsRRtnHv3hvosa7pndWgZwfOBYpmAybLh6mfgjADUXxtvBzPMT7NGab2rv5ORiv8y4FvOQ45xeKwNKz0Wr2wxiD4tfyzop3_D9OB-ta3F6E') 
	RETURNING *;`

	user := new(entity.User)

	if err := u.db.Raw(query).Scan(user).Error; err != nil {
		u.logger.WithFields(logrus.Fields{
			"location": "util.UserTest/Create",
			"section":  "gorm.DB.Raw",
		}).Errorf(err.Error())
	}

	return user
}

func (u *UserTest) Delete() {
	if err := u.db.Exec("DELETE FROM users;").Error; err != nil {
		u.logger.WithFields(logrus.Fields{
			"location": "util.UserTest/Delete",
			"section":  "orm.DB.Exec",
		}).Errorf(err.Error())
	}
}