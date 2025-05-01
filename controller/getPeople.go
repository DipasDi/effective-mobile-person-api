package controller

import (
	"EffectiveMobile/database"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// @Summary Получить список людей
// @Description Возвращает список людей с пагинацией и фильтрацией
// @Tags persons
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит на страницу" default(10)
// @Param name query string false "Фильтр по имени"
// @Param surname query string false "Фильтр по фамилии"
// @Param age query int false "Фильтр по возрасту"
// @Param gender query string false "Фильтр по полу"
// @Param nationality query string false "Фильтр по национальности"
// @Success 200 {object} map[string]interface{}
// @Router /api/getPeople [get]
func GetPeople(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	name := c.Query("name")
	surname := c.Query("surname")
	age := c.Query("age")
	gender := c.Query("gender")
	nationality := c.Query("nationality")

	database.Dbconnect()
	defer database.Db.Close()

	baseQuery := "SELECT id, name, surname, patronymic, age, gender, nationality FROM people"
	countQuery := "SELECT COUNT(*) FROM people"

	var conditions []string
	var args []interface{}
	var countArgs []interface{}

	if name != "" {
		conditions = append(conditions, "name LIKE $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, "%"+name+"%")
		countArgs = append(countArgs, "%"+name+"%")
	}

	if surname != "" {
		conditions = append(conditions, "surname LIKE $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, "%"+surname+"%")
		countArgs = append(countArgs, "%"+surname+"%")
	}

	if age != "" {
		conditions = append(conditions, "age = $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, age)
		countArgs = append(countArgs, age)
	}

	if gender != "" {
		conditions = append(conditions, "gender = $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, gender)
		countArgs = append(countArgs, gender)
	}

	if nationality != "" {
		conditions = append(conditions, "nationality = $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, nationality)
		countArgs = append(countArgs, nationality)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	dataQuery := baseQuery + whereClause +
		fmt.Sprintf(" ORDER BY id LIMIT %d OFFSET %d", limit, offset)

	rows, err := database.Db.Query(dataQuery, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query people")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query people"})
		return
	}
	defer rows.Close()

	var people []EnrichedPerson
	for rows.Next() {
		var p EnrichedPerson
		err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan person")
			continue
		}
		people = append(people, p)
	}

	fullCountQuery := countQuery + whereClause
	var total int
	err = database.Db.QueryRow(fullCountQuery, countArgs...).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count people")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count people"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": people,
		"meta": gin.H{
			"total":  total,
			"page":   page,
			"limit":  limit,
			"pages":  int(math.Ceil(float64(total) / float64(limit))),
			"offset": offset,
		},
	})
}
