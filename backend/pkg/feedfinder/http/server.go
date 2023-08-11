package http

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/spezifisch/rueder3/backend/docs" // api docs generated by Swag CLI
	"github.com/spezifisch/rueder3/backend/internal/common"
	"github.com/spezifisch/rueder3/backend/pkg/feedfinder/controller"
)

// Server is a http server
type Server struct {
	Bind string

	engine            *gin.Engine
	controller        *controller.Controller
	jwtSecretKey      string
	isDevelopmentMode bool
	trustedProxies    []string
}

// NewServer creates a default http backend
func NewServer(controller *controller.Controller, bind string, jwtSecretKey string, isDevelopmentMode bool, trustedProxies []string) *Server {
	s := &Server{
		Bind:              bind,
		controller:        controller,
		jwtSecretKey:      jwtSecretKey,
		isDevelopmentMode: isDevelopmentMode,
		trustedProxies:    trustedProxies,
	}
	s.init()
	return s
}

func (s *Server) init() {
	if !s.isDevelopmentMode {
		gin.SetMode("release")
	}

	s.engine = gin.Default()
	common.GinSetTrustedProxies(s.engine, s.trustedProxies)

	if s.isDevelopmentMode {
		s.engine.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))
	}

	s.initAPI()
}

// Run starts the server
func (s *Server) Run() {
	err := s.engine.Run(s.Bind)
	if err != nil {
		log.WithError(err).Error("gin failed")
	}
}
