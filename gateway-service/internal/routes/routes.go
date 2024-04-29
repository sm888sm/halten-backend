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
		boardRoutes.GET("/:boardID", boardHandler.GetBoardByID)
		boardRoutes.GET("/", boardHandler.GetBoardList)
		boardRoutes.PUT("/:boardID/name", boardHandler.UpdateBoardName)
		boardRoutes.DELETE("/:boardID", boardHandler.DeleteBoard)
		boardRoutes.POST("/:boardID/users", boardHandler.AddBoardUsers)
		boardRoutes.DELETE("/:boardID/users", boardHandler.RemoveBoardUsers)
		boardRoutes.GET("/:boardID/users", boardHandler.GetBoardMembers)
		boardRoutes.PUT("/:boardID/users/:userID/role", boardHandler.AssignBoardUsersRole)
		boardRoutes.PUT("/:boardID/owner", boardHandler.ChangeBoardOwner)
		boardRoutes.PUT("/:boardID/visibility", boardHandler.ChangeBoardVisibility)
		boardRoutes.POST("/:boardID/labels", boardHandler.AddLabel)
		boardRoutes.DELETE("/:boardID/labels/:labelID", boardHandler.RemoveLabel)
		boardRoutes.PUT("/:boardID/archive", boardHandler.ArchiveBoard)
		boardRoutes.PUT("/:boardID/restore", boardHandler.RestoreBoard)
		boardRoutes.GET("/archived", boardHandler.GetArchivedBoardList)
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
		cardRoutes.GET("/board/:board-id", cardHandler.GetCardsByBoard)
		cardRoutes.GET("/list/:list-id", cardHandler.GetCardsByList)
		cardRoutes.GET("/:card-id", cardHandler.GetCardByID)
		cardRoutes.DELETE("/:card-id", cardHandler.DeleteCard)
		cardRoutes.PUT("/:card-id/position", cardHandler.MoveCardPosition)
	}
}
