package routes

import (
	"github.com/gin-gonic/gin"
	external_services "github.com/sm888sm/halten-backend/gateway-service/external/services"
	"github.com/sm888sm/halten-backend/gateway-service/internal/handlers"
	"github.com/sm888sm/halten-backend/gateway-service/internal/middlewares"
)

func SetupRoutes(r *gin.Engine, svc *external_services.Services, secretKey string) {

	userHandler := handlers.NewUserHandler(svc)
	authHandler := handlers.NewAuthHandler(svc)
	boardHandler := handlers.NewBoardHandler(svc)
	listHandler := handlers.NewListHandler(svc)
	cardHandler := handlers.NewCardHandler(svc)

	userRoutes := r.Group("/user")
	userRoutes.POST("/create", userHandler.CreateUser)
	userRoutes.PUT("/confirm-new-email", userHandler.ConfirmNewEmail)

	userRoutes.Use(middlewares.UserMiddleware(svc, secretKey))
	{
		userRoutes.PUT("/update-email", userHandler.UpdateEmail)
		userRoutes.PUT("/update-password", userHandler.UpdatePassword)
		userRoutes.PUT("/update-username", userHandler.UpdateUsername)
	}

	authRoutes := r.Group("/auth")
	authRoutes.POST("/login", authHandler.Login)
	authRoutes.POST("/refresh", authHandler.RefreshToken)

	boardRoutes := r.Group("/boards")
	boardRoutes.Use(middlewares.UserMiddleware(svc, secretKey))
	{
		boardRoutes.POST("/", boardHandler.CreateBoard)
		boardRoutes.GET("/:id", boardHandler.GetBoardByID)
		boardRoutes.GET("/", boardHandler.GetBoardList)
		boardRoutes.PUT("/:id", boardHandler.UpdateBoard)
		boardRoutes.DELETE("/:id", boardHandler.DeleteBoard)
		boardRoutes.POST("/:id/users", boardHandler.AddBoardUsers)
		boardRoutes.DELETE("/:id/users", boardHandler.RemoveBoardUsers)
		boardRoutes.GET("/:id/users", boardHandler.GetBoardUsers)
		boardRoutes.PUT("/:id/users/:userId", boardHandler.AssignBoardUserRole)
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
