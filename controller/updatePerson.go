package controller

import (
	"EffectiveMobile/database"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary Обновить данные человека
// @Description Обновляет данные человека по ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "ID человека"
// @Param input body UpdatePersonInput true "Обновленные данные"
// @Success 200
// @Failure 400 {object} map[string]string
// @Router /api/update/{id} [put]
func UpdatePerson(c *gin.Context) {
	id := c.Param("id")

	var input UpdatePersonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Msg("Invalid input data")
		return
	}

	query := "UPDATE people SET "
	var args []interface{}
	var setClauses []string
	argPos := 1

	if *input.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf(" name = $%d", argPos))
		args = append(args, *input.Name)
		argPos++
	}
	if *input.Surname != "" {
		setClauses = append(setClauses, fmt.Sprintf(" surname = $%d", argPos))
		args = append(args, *input.Surname)
		argPos++
	}
	if *input.Patronymic != "" {
		setClauses = append(setClauses, fmt.Sprintf(" patronymic = $%d", argPos))
		args = append(args, *input.Patronymic)
		argPos++
	}
	if *input.Age != 0 {
		setClauses = append(setClauses, fmt.Sprintf(" age = $%d", argPos))
		args = append(args, *input.Age)
		argPos++
	}
	if *input.Gender != "" {
		setClauses = append(setClauses, fmt.Sprintf(" gender = $%d", argPos))
		args = append(args, *input.Gender)
		argPos++
	}
	if *input.Nationality != "" {
		setClauses = append(setClauses, fmt.Sprintf(" nationality = $%d", argPos))
		args = append(args, *input.Nationality)
		argPos++
	}

	query += strings.Join(setClauses, ",")
	query += fmt.Sprintf(" WHERE id = $%d", argPos)
	args = append(args, id)

	database.Dbconnect()
	defer database.Db.Close()

	if _, err := database.Db.Exec(query, args...); err != nil {
		log.Error().Err(err).Str("query", query).Msg("Failed to update person")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person"})
		return
	}
}
