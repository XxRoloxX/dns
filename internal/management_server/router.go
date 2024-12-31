package managementserver

import "github.com/gin-gonic/gin"

type Controller interface {
	Handle(ctx *gin.Context)
}

type RecordsRouterParams struct {
	Engine                 *gin.Engine
	NewRecordController    Controller
	DeleteRecordController Controller
	GetRecordsController   Controller
}

func NewRecordsRouter(params *RecordsRouterParams) *gin.RouterGroup {

	router := params.Engine.Group("/records")

	router.POST("", params.NewRecordController.Handle)
	router.GET("", params.GetRecordsController.Handle)
	router.DELETE("/:id", params.DeleteRecordController.Handle)

	return router
}
