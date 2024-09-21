package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Users struct {
	Id        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func main() {
	r := gin.Default()
	godotenv.Load()

	v1 := r.Group("api/v1")
	{
		v1.POST("/users", PostUser)
		v1.GET("/users", GetUsers)
		v1.GET("/users/:id", GetUser)
		v1.PUT("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)
	}

	r.Run(":8080")
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func InitDb() *sql.DB {
	postgresUrl := goDotEnvVariable("NEON_POSTGRES_URL")
	connStr := postgresUrl
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	return db
}

func PostUser(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var user Users
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if user.Firstname != "" && user.Lastname != "" {
		query := `INSERT INTO users (firstname, lastname) VALUES ($1, $2)`
		_, err := db.Exec(query, user.Firstname, user.Lastname)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"success": user})
	} else {
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}
}

func GetUsers(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	rows, err := db.Query("SELECT id, firstname, lastname FROM users")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []Users
	for rows.Next() {
		var user Users
		if err := rows.Scan(&user.Id, &user.Firstname, &user.Lastname); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	user_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	db := InitDb()
	defer db.Close()

	var user Users
	err = db.QueryRow("SELECT id, firstname, lastname FROM users WHERE id = $1", user_id).Scan(&user.Id, &user.Firstname, &user.Lastname)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "User not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {
	// Future code...
}

func DeleteUser(c *gin.Context) {
	// Future code...
}
