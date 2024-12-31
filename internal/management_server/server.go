package managementserver

import "github.com/gin-gonic/gin"

type Server struct {
	engine *gin.Engine
}

func NewServer() *Server {
	return &Server{
		engine: gin.Default(),
	}
}

func (s *Server) Start() {

	repository := NewPostgresRecordsRepository()
	service := NewRecordsService(repository)

	_ = NewRecordsRouter(&RecordsRouterParams{
		Engine:                 s.engine,
		NewRecordController:    NewNewRecordController(service),
		GetRecordsController:   NewGetRecordsController(service),
		DeleteRecordController: NewDeleteRecordController(service),
	})

	s.engine.Run(":8080")
}
