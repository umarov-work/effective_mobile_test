package handlers

import (
	"effective_mobile_test/internal/models"
	"effective_mobile_test/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// CreatePerson
// @Summary Create a new person
// @Description Creates a new person with enrichment from external APIs
// @Tags persons
// @Accept json
// @Produce json
// @Param person body models.PersonInput true "Person data"
// @Success 201 {object} models.PersonResponse
// @Failure 400 {object} map[string]string "Invalid JSON"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /person [post]
func CreatePerson(c *gin.Context, db *gorm.DB, logger *logrus.Logger) {
	var personInput models.PersonInput
	err := c.BindJSON(&personInput)
	if err != nil {
		logger.Error("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	person := &models.Person{
		Name:       personInput.Name,
		Surname:    personInput.Surname,
		Patronymic: personInput.Patronymic,
	}

	err = services.EnrichPersonWithAge(person)
	if err != nil {
		logger.Error("Failed to enrich person with age:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Enrichment with age failed"})
		return

	}

	err = services.EnrichPersonWithGender(person)
	if err != nil {
		logger.Error("Failed to enrich person with gender:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Enrichment with gender failed"})
		return
	}

	err = services.EnrichPersonWithNationality(person)
	if err != nil {
		logger.Error("Failed to enrich person with nationality:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Enrichment with nationality failed"})
		return
	}

	err = db.Create(person).Error
	if err != nil {
		logger.Error("Failed to create person:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create person failed"})
		return
	}

	response := &models.PersonResponse{
		ID:          person.ID,
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Age:         person.Age,
		Gender:      person.Gender,
		Nationality: person.Nationality,
	}

	logger.Info("Person created:", person.ID)
	c.JSON(http.StatusCreated, response)
}

// DeletePerson
// @Summary Delete a person
// @Description Deletes a person by ID
// @Tags persons
// @Produce json
// @Param id path int true "Person ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string "Person not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /person/{id} [delete]
func DeletePerson(c *gin.Context, db *gorm.DB, logger *logrus.Logger) {
	id := c.Param("id")
	result := db.Delete(&models.Person{}, id)
	if result.Error != nil {
		logger.Error("Failed to delete person:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete person"})
		return
	}
	if result.RowsAffected == 0 {
		logger.Warn("Person not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	logger.Info("Person deleted:", id)
	c.Status(http.StatusNoContent)
}

// GetPersons
// @Summary Get list of persons
// @Description Retrieve a list of persons with optional filters and pagination
// @Tags persons
// @Produce json
// @Param name query string false "Filter by name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} models.PersonResponse
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /persons [get]
func GetPersons(c *gin.Context, db *gorm.DB, logger *logrus.Logger) {
	var persons []models.Person
	query := db

	if name := c.Query("name"); name != "" {
		query = query.Where("name = ?", name)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	err := query.Limit(limit).Offset(offset).Find(&persons).Error
	if err != nil {
		logger.Error("Failed to fetch persons:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Преобразуем []Person в []PersonResponse
	var personResponses []models.PersonResponse
	for _, person := range persons {
		personResponses = append(personResponses, models.PersonResponse{
			ID:          person.ID,
			Name:        person.Name,
			Surname:     person.Surname,
			Patronymic:  person.Patronymic,
			Age:         person.Age,
			Gender:      person.Gender,
			Nationality: person.Nationality,
		})
	}

	logger.Info("Fetched persons, count:", len(personResponses))
	c.JSON(http.StatusOK, personResponses)
}

// UpdatePerson
// @Summary Update a person
// @Description Updates an existing person by ID with enrichment
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body models.PersonInput true "Person data"
// @Success 200 {object} models.PersonResponse
// @Failure 400 {object} map[string]string "Invalid JSON"
// @Failure 404 {object} map[string]string "Person not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /person/{id} [put]
func UpdatePerson(c *gin.Context, db *gorm.DB, logger *logrus.Logger) {
	id := c.Param("id")
	var person models.Person

	err := db.First(&person, id).Error
	if err != nil {
		logger.Warn("Person not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	var personInput models.PersonInput
	err = c.BindJSON(&personInput)
	if err != nil {
		logger.Error("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	person.Name = personInput.Name
	person.Surname = personInput.Surname
	person.Patronymic = personInput.Patronymic

	err = services.EnrichPersonWithAge(&person)
	if err != nil {
		logger.Error("Failed to enrich person with age:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Enrichment with age failed"})
		return
	}

	err = services.EnrichPersonWithGender(&person)
	if err != nil {
		logger.Error("Failed to enrich person with gender:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Enrichment with gender failed"})
		return
	}

	err = services.EnrichPersonWithNationality(&person)
	if err != nil {
		logger.Error("Failed to enrich person with nationality:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Enrichment with nationality failed"})
		return
	}

	err = db.Save(&person).Error
	if err != nil {
		logger.Error("Failed to update person:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	response := &models.PersonResponse{
		ID:          person.ID,
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Age:         person.Age,
		Gender:      person.Gender,
		Nationality: person.Nationality,
	}

	logger.Info("Person updated:", id)
	c.JSON(http.StatusOK, response)
}
