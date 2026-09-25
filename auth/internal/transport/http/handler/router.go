package handler

import "github.com/gin-gonic/gin"

func HTTPRouter(
	h *HTTPHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.RegisterUser)
		auth.POST("/login", h.LoginUser)
		auth.POST("/refresh", h.RefreshToken)

		protected := auth.Group("/")
		protected.Use(h.RequireAuth())
		{
			protected.GET("/user", h.GetUser)
			protected.PATCH("/user", h.UpdateUser)
			protected.DELETE("/user", h.DeleteUser)
		}
	}

	return r
}
