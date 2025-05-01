// @title Effective Mobile API
// @version 1.0
// @description API для работы с данными людей
// @host localhost:8080
// @BasePath /
package main

import (
	"EffectiveMobile/controller"
	"EffectiveMobile/database"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()

	if err := database.RunMigrations(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
		return
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Rest Запросы
	r.GET("/api/getPeople", controller.GetPeople)
	r.POST("/api/enrich", controller.EnrichPersonHandler)
	r.DELETE("/api/people/:id", controller.DeletePerson)
	r.PUT("/api/update/:id", controller.UpdatePerson)

	log.Info().Msg("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
