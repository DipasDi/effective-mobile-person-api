package controller

import (
	"EffectiveMobile/database"
	"EffectiveMobile/external"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary Обогатить данные человека
// @Description Добавляет нового человека с обогащенными данными (возраст, пол, национальность)
// @Tags persons
// @Accept json
// @Produce json
// @Param input body EnrichedPerson true "Данные человека"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/enrich [post]

func EnrichPersonHandler(c *gin.Context) {
	var input EnrichedPerson
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Msg("Invalid input data")
		return
	}

	age, err := external.GetAge(input.Name)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get age")
		return
	}

	gender, err := external.GetGender(input.Name)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get gender")
		return
	}

	nationality, err := external.GetNationality(input.Name)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get nationality")
		return
	}

	database.Dbconnect()
	defer database.Db.Close()

	result, err := database.Db.Exec("INSERT INTO people (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6)", input.Name, input.Surname, input.Patronymic, *age, *gender, *nationality)
	if err != nil {
		log.Info().Err(err).Msg("Failed to insert person")
		return
	} else {
		log.Info().Msg("Successfully enriched person data")
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Person data successfully enriched and saved",
		"data": gin.H{
			"name":        input.Name,
			"surname":     input.Surname,
			"patronymic":  input.Patronymic,
			"age":         *age,
			"gender":      *gender,
			"nationality": *nationality,
			"result":      result,
		},
	})
}
