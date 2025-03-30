package main

import (
	"effective_mobile_test/config"
	"effective_mobile_test/internal/database"
	"effective_mobile_test/internal/handlers"
	"effective_mobile_test/internal/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	log := logger.InitLogger()
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
		log.Info("Exiting...")
		return
	}
	log.Debug("Loaded config:", conf)

	db, err := database.ConnectDB(conf.DbHost, conf.DbPort, conf.DbUser, conf.DbPassword, conf.DbName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		log.Info("Exiting...")
		return
	}
	log.Debug("Connected to DB:", db)

	err = database.Migrate(db)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Info("Migration complete")

	router := gin.Default()
	router.POST("/person", func(c *gin.Context) {
		handlers.CreatePerson(c, db, log)
	})
	router.GET("/persons", func(c *gin.Context) {
		handlers.GetPersons(c, db, log)
	})
	router.PUT("/person", func(c *gin.Context) {
		handlers.UpdatePerson(c, db, log)
	})
	router.DELETE("/person", func(c *gin.Context) {
		handlers.DeletePerson(c, db, log)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))
	router.StaticFile("/docs/swagger.json", "../docs/swagger.json")
	port := conf.Port
	log.Info("Starting server on port:", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
		log.Info("Exiting...")
		return
	}
	log.Info("Service started on port:", port)
}
