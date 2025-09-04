package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"image_server/internal/global"
)

func LogMiddleware(c *gin.Context) {
	log := global.Log
	uid := uuid.New().String()
	logger := log.WithField("logID", uid)
	c.Set("log", logger)
}

func GetLog(c *gin.Context) *logrus.Entry {
	return c.MustGet("log").(*logrus.Entry)
}
