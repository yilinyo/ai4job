package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yilin/ai-for-backend/internal/interface/handler"
)

// Dependencies holds all handler and middleware dependencies for the router.
type Dependencies struct {
	UserHandler          *handler.UserHandler
	ResumeRepoHandler    *handler.ResumeRepoHandler
	ResumeVersionHandler *handler.ResumeVersionHandler
	AuthMiddleware       gin.HandlerFunc
}

// NewRouter creates and configures the Gin engine with all routes.
func NewRouter(deps Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")

	user := v1.Group("/user")
	{
		user.POST("/email-code", deps.UserHandler.SendEmailCode)
		user.POST("/register", deps.UserHandler.Register)
		user.POST("/login", deps.UserHandler.Login)

		authed := user.Group("/", deps.AuthMiddleware)
		{
			authed.GET("/profile", deps.UserHandler.GetProfile)
			authed.PUT("/profile", deps.UserHandler.UpdateProfile)
			authed.POST("/logout", deps.UserHandler.Logout)
		}
	}

	auth := v1.Group("/", deps.AuthMiddleware)
	{
		auth.GET("/resume-repos", deps.ResumeRepoHandler.List)
		auth.POST("/resume-repos", deps.ResumeRepoHandler.Create)

		auth.GET("/resume-repos/:id/versions", deps.ResumeVersionHandler.List)
		auth.POST("/resume-repos/:id/versions", deps.ResumeVersionHandler.Create)
		auth.GET("/resume-repos/:id/versions/:id", deps.ResumeVersionHandler.Get)
		auth.PUT("/resume-repos/:id/versions/:id", deps.ResumeVersionHandler.Update)
		auth.DELETE("/resume-repos/:id/versions/:id", deps.ResumeVersionHandler.Delete)
		auth.POST("/resume-repos/:id/versions/:id/set-default", deps.ResumeVersionHandler.SetDefault)

		auth.GET("/resume-repos/:id", deps.ResumeRepoHandler.Get)
		auth.PUT("/resume-repos/:id", deps.ResumeRepoHandler.Update)
		auth.DELETE("/resume-repos/:id", deps.ResumeRepoHandler.Delete)
	}

	return r
}
