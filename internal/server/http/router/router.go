package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"quiz-service/init/config"
	"quiz-service/init/logger"
	"quiz-service/internal/repository"
	"quiz-service/internal/server/http/handlers"
	"quiz-service/internal/server/http/middleware"
	"quiz-service/internal/service"
	"quiz-service/pkg/hash"
)

type Router struct {
	router  *gin.RouterGroup
	handler *handlers.Handler
}

func InitRouterAndComponents(router *gin.RouterGroup, db *sqlx.DB, cfg *config.Config, httpLogger, dbLogger, quizLogger *logger.Logger) *Router {
	repo := repository.NewRepository(db, dbLogger)
	hasher := hash.NewSHA512Hasher(cfg.PasswordSalt)
	serv := service.NewService(repo, hasher, quizLogger)
	handler := handlers.NewHandler(serv, httpLogger)

	return &Router{
		router:  router,
		handler: handler,
	}
}

func (r *Router) Routes() {
	r.router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", nil)
	})

	r.router.POST("/register", r.handler.Register)
	r.router.POST("/login", r.handler.Login)

	user := r.router.Group("/:userId", middleware.UserId(), r.handler.Authenticated)
	{
		user.POST("/quit", r.handler.Quit)

		variants := user.Group("/variant")
		{
			variants.GET("/", func(ctx *gin.Context) {
				ctx.HTML(http.StatusOK, "variants.html", nil)
			})

			variants.POST("/add", r.handler.VariantAdd)
			variants.GET("/list", r.handler.VariantList)

			variantName := variants.Group("/:variantName", r.handler.VariantCheck)
			{
				variantName.GET("/", func(ctx *gin.Context) {
					ctx.HTML(http.StatusOK, "test.html", nil)
				})

				variantName.DELETE("/remove", r.handler.VariantRemove)
				variantName.POST("/start", r.handler.VariantStart)
				variantName.GET("/results", r.handler.VariantResults)
				variantName.GET("/get", r.handler.VariantGet)

				question := variantName.Group("/question")
				{
					question.POST("/add", r.handler.QuestionAdd)
					question.DELETE("/remove", r.handler.QuestionRemove)

					questionId := question.Group("/:questionId", middleware.QuestionId())
					{
						questionId.GET("/get", middleware.QuestionId(), r.handler.QuestionGet)
						questionId.POST("/accept", middleware.QuestionId(), r.handler.QuestionAccept)
					}
				}
			}
		}
	}
}
