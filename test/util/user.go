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
	query := "INSERT INTO users(email, full_name, password) VALUES('johndoe@gmail.com', 'John Doe', 'rahasia') RETURNING *;"

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