package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sm888sm/halten-backend/gateway-service/internal/handlers"
	"github.com/sm888sm/halten-backend/gateway-service/internal/middlewares"
	"github.com/sm888sm/halten-backend/gateway-service/internal/services"
)

func SetupRoutes(r *gin.Engine, svc *services.Services, secretKey string) {

	authHandler := handlers.NewAuthHandler(svc)
	userHandler := handlers.NewUserHandler(svc)
	boardHandler := handlers.NewBoardHandler(svc)
	listHandler := handlers.NewListHandler(svc)
	cardHandler := handlers.NewCardHandler(svc)

	authRoutes := r.Group("/auth")
	authRoutes.POST("/login", authHandler.Login)
	authRoutes.POST("/refresh", authHandler.RefreshToken)

	userRoutes := r.Group("/user")
	userRoutes.POST("/create", userHandler.CreateUser)
	userRoutes.PUT("/confirm-new-email", userHandler.ConfirmNewEmail)

	userRoutes.Use(middlewares.UserMiddleware(svc, secretKey))
	{
		userRoutes.PUT("/update-email", userHandler.UpdateEmail)
		userRoutes.PUT("/update-password", userHandler.UpdatePassword)
		userRoutes.PUT("/update-username", userHandler.UpdateUsername)
	}

	boardRoutes := r.Group("/boards")
	boardRoutes.Use(middlewares.UserMiddleware(svc, secretKey))
	{
		boardRoutes.POST("/", boardHandler.CreateBoard)
		boardRoutes.GET("/:id", boardHandler.GetBoardByID)
		boardRoutes.GET("/", boardHandler.GetBoards)
		boardRoutes.PUT("/:id", boardHandler.UpdateBoard)
		boardRoutes.DELETE("/:id", boardHandler.DeleteBoard)
		boardRoutes.POST("/:id/users", boardHandler.AddBoardUsers)
		boardRoutes.DELETE("/:id/users", boardHandler.RemoveBoardUsers)
		boardRoutes.GET("/:id/users", boardHandler.GetBoardUsers)
		boardRoutes.PUT("/:id/users/:userId", boardHandler.AssignUserRoleBoard)
		boardRoutes.PUT("/:id/owner", boardHandler.ChangeBoardOwner)
	}

	listRoutes := r.Group("/lists")
	listRoutes.Use(middlewares.UserMiddleware(svc, secretKey))
	{
		listRoutes.POST("/", listHandler.CreateList)
		listRoutes.GET("/:id", listHandler.GetListsByBoard)
		listRoutes.PUT("/:id", listHandler.UpdateList)
		listRoutes.DELETE("/:id", listHandler.DeleteList)
		listRoutes.PUT("/:id/position", listHandler.MoveListPosition)
	}

	cardRoutes := r.Group("/cards")
	cardRoutes.Use(middlewares.UserMiddleware(svc, secretKey))
	{
		cardRoutes.POST("/", cardHandler.CreateCard)
		cardRoutes.GET("/:id", cardHandler.GetCardsByList)
		cardRoutes.PUT("/:id", cardHandler.UpdateCard)
		cardRoutes.DELETE("/:id", cardHandler.DeleteCard)
		cardRoutes.PUT("/:id/position", cardHandler.MoveCardPosition)
	}
}
