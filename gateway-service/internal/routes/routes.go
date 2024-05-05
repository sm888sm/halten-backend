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
	userRoutes.PUT("/confirm-email", userHandler.ConfirmEmail)

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
		boardRoutes.GET("/", boardHandler.GetBoardList)
		boardRoutes.GET("/:boardID", boardHandler.GetBoardByID)
		boardRoutes.GET("/:boardID/users", boardHandler.GetBoardMembers)
		boardRoutes.GET("/archived", boardHandler.GetArchivedBoardList)

		boardRoutes.POST("/", boardHandler.CreateBoard)
		boardRoutes.POST("/:boardID/users", boardHandler.AddBoardUsers)
		boardRoutes.POST("/:boardID/labels", boardHandler.AddLabel)

		boardRoutes.PUT("/:boardID/name", boardHandler.UpdateBoardName)
		boardRoutes.PUT("/:boardID/users/:userID/role", boardHandler.AssignBoardUsersRole)
		boardRoutes.PUT("/:boardID/owner", boardHandler.ChangeBoardOwner)
		boardRoutes.PUT("/:boardID/visibility", boardHandler.ChangeBoardVisibility)
		boardRoutes.PUT("/:boardID/archive", boardHandler.ArchiveBoard)
		boardRoutes.PUT("/:boardID/restore", boardHandler.RestoreBoard)

		boardRoutes.DELETE("/:boardID", boardHandler.DeleteBoard)
		boardRoutes.DELETE("/:boardID/users", boardHandler.RemoveBoardUsers)
		boardRoutes.DELETE("/:boardID/labels/:labelID", boardHandler.RemoveLabel)

	}

	listRoutes := r.Group("/lists")
	listRoutes.Use(middlewares.UserMiddleware(svc, secretKey))
	{
		listRoutes.GET("/:listID", listHandler.GetListByID)
		listRoutes.GET("/board/:boardID", listHandler.GetListsByBoard)

		listRoutes.POST("/", listHandler.CreateList)

		listRoutes.PUT("/:listID", listHandler.UpdateListName)
		listRoutes.PUT("/:listID/move", listHandler.MoveListPosition)
		listRoutes.PUT("/:listID/archive", listHandler.ArchiveList)
		listRoutes.PUT("/:listID/restore", listHandler.RestoreList)

		listRoutes.DELETE("/:listID", listHandler.DeleteList)
	}

	cardRoutes := r.Group("/cards")
	cardRoutes.Use(middlewares.UserMiddleware(svc, secretKey))
	{
		cardRoutes.GET("/:cardID", cardHandler.GetCardByID)
		cardRoutes.GET("/list/:listID", cardHandler.GetCardsByList)
		cardRoutes.GET("/board/:boardID", cardHandler.GetCardsByBoard)

		cardRoutes.POST("/", cardHandler.CreateCard)
		cardRoutes.POST("/:cardID/attachment/:attachmentID", cardHandler.AddCardAttachment)
		cardRoutes.POST("/:cardID/comment", cardHandler.AddCardComment)

		cardRoutes.PUT("/:cardID/position", cardHandler.MoveCardPosition)
		cardRoutes.PUT("/:cardID/name", cardHandler.UpdateCardName)
		cardRoutes.PUT("/:cardID/label/:labelID", cardHandler.AddCardLabel)
		cardRoutes.PUT("/:cardID/dates", cardHandler.SetCardDates)
		cardRoutes.PUT("/:cardID/completed", cardHandler.ToggleCardCompleted)
		cardRoutes.PUT("/:cardID/members", cardHandler.AddCardMembers)
		cardRoutes.PUT("/:cardID/archive", cardHandler.ArchiveCard)
		cardRoutes.PUT("/:cardID/restore", cardHandler.RestoreCard)

		cardRoutes.DELETE("/:cardID/label/:labelID", cardHandler.RemoveCardLabel)
		cardRoutes.DELETE("/:cardID/attachment/:attachmentID", cardHandler.RemoveCardAttachment)
		cardRoutes.DELETE("/:cardID/comment/:commentID", cardHandler.RemoveCardComment)
		cardRoutes.DELETE("/:cardID/members", cardHandler.RemoveCardMembers)
		cardRoutes.DELETE("/:cardID", cardHandler.DeleteCard)
	}
}
