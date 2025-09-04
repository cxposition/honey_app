package user_service

import (
	"github.com/sirupsen/logrus"
)

type UserService struct {
	log *logrus.Entry
}

func NewUserService(log *logrus.Entry) *UserService {
	return &UserService{
		log: log,
	}
}
