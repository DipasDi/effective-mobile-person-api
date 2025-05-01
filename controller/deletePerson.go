package controller

import (
	"EffectiveMobile/database"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary Удалить данные человека
// @Description Удаляет данные человека по ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "ID человека"
// @Success 200
// @Failure 400 {object} map[string]string
// @Router /api/people/{id} [delete]
func DeletePerson(c *gin.Context) {
	id := c.Param("id")

	database.Dbconnect()
	defer database.Db.Close()

	if _, err := database.Db.Exec("DELETE FROM people WHERE id = $1", id); err != nil {
		log.Info().Err(err).Msg("Failed to delete person")
	}
}
