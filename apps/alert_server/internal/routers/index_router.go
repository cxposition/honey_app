package routers

import (
	"alert_server/internal/api"
	"github.com/gin-gonic/gin"
)

func IndexRouter(r *gin.RouterGroup) {
	app := api.App.IndexApi
	r.GET("signature_agg", app.SignatureAggView)
}
