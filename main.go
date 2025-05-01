package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"hello-world-devops-backend/models"
)

// --- Global Database Variable ---
var DB *gorm.DB

// LoadEnv loads environment variables from a single `.env` file in the root.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v. Relying on system environment variables.", err)
	} else {
		log.Println("Loaded environment variables from .env")
	}
}

// ConnectDatabase initializes the database connection using environment variables.
func ConnectDatabase() {
	LoadEnv()
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	portStr := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// Set defaults if variables are missing
	if host == "" {
		host = "localhost"
	}
	if portStr == "" {
		portStr = "5432"
	}
	if sslmode == "" {
		sslmode = "disable"
	}
	if user == "" || dbname == "" {
		log.Fatal("Essential Database environment variables (DB_USER, DB_NAME) are not set.")
	}

	// Parse port
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT '%s': %v", portStr, err)
	}

	// Construct the Data Source Name (DSN) for PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode)

	// Open the database connection
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Successfully connected to database: %s", dbname)

	DB = database
}

// --- Handlers ---

// CreateTask handles POST requests to create a new task.
func CreateTask(c *gin.Context) {
	var input models.CreateTaskInput
	// Validate input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new task object
	task := models.Task{Title: input.Title, Description: input.Description, Completed: false} // Use models.Task

	// Save the task to the database
	result := DB.Create(&task)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": result.Error.Error()})
		return
	}

	// Return the created task
	c.JSON(http.StatusCreated, task)
}

// FindTasks handles GET requests to retrieve all tasks.
func FindTasks(c *gin.Context) {
	var tasks []models.Task
	// Find all tasks in the database
	result := DB.Find(&tasks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks", "details": result.Error.Error()})
		return
	}

	// Return the list of tasks
	c.JSON(http.StatusOK, tasks)
}

// FindTask handles GET requests to retrieve a single task by ID.
func FindTask(c *gin.Context) {
	var task models.Task
	taskID := c.Param("id")

	// Find the task by ID
	result := DB.First(&task, taskID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task", "details": result.Error.Error()})
		}
		return
	}

	// Return the found task
	c.JSON(http.StatusOK, task)
}

// UpdateTask handles PUT requests to update an existing task by ID.
func UpdateTask(c *gin.Context) {
	taskID := c.Param("id")

	// Find the existing task
	var task models.Task
	findResult := DB.First(&task, taskID)
	if findResult.Error != nil {
		if findResult.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find task for update", "details": findResult.Error.Error()})
		}
		return
	}

	// Validate input JSON
	var input models.UpdateTaskInput // Use models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a map to hold only the fields to be updated
	updateData := make(map[string]interface{})
	if input.Title != nil {
		updateData["title"] = *input.Title
	}
	if input.Description != nil {
		updateData["description"] = *input.Description
	}
	if input.Completed != nil {
		updateData["completed"] = *input.Completed
	}

	// Check if there's anything to update
	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No update fields provided"})
		return
	}

	// Update the task attributes in the database using the map
	updateResult := DB.Model(&task).Updates(updateData)
	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task", "details": updateResult.Error.Error()})
		return
	}

	// Return the updated task
	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE requests to remove a task by ID.
func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	// Delete the task by ID
	result := DB.Delete(&models.Task{}, taskID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task", "details": result.Error.Error()})
		return
	}

	// Check if any row was actually affected
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or already deleted"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// --- Main Function ---

func main() {
	// Connect to the database
	ConnectDatabase()

	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)
	log.Printf("Gin mode set to: %s", ginMode)

	// Initialize Gin router
	router := gin.Default()

	// --- CORS Middleware ---
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" && ginMode == gin.DebugMode {
		frontendOrigin = "http://localhost:5173"
		log.Printf("Warning: FRONTEND_ORIGIN not set, allowing default %s for CORS", frontendOrigin)
	}

	corsConfig := cors.DefaultConfig()
	if frontendOrigin != "" {
		corsConfig.AllowOrigins = []string{frontendOrigin}
	} else {
		corsConfig.AllowAllOrigins = true
		log.Println("Warning: No FRONTEND_ORIGIN set, allowing all origins for CORS.")
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true

	router.Use(cors.New(corsConfig))

	// --- API Routes ---
	api := router.Group("/api/v1")
	{
		api.POST("/tasks", CreateTask)
		api.GET("/tasks", FindTasks)
		api.GET("/tasks/:id", FindTask)
		api.PUT("/tasks/:id", UpdateTask)
		api.DELETE("/tasks/:id", DeleteTask)
	}

	// Simple health check route
	router.GET("/health", func(c *gin.Context) {
		sqlDB, err := DB.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "DOWN", "error": "Failed to get DB connection"})
			return
		}
		err = sqlDB.Ping()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "DOWN", "error": "Failed to ping DB"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// --- Start Server ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
