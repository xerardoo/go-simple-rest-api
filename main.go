package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Mysql Library : https://github.com/go-sql-driver/mysql
// https://idineshkrishnan.com/crud-operations-with-mysql-in-go-language/

// variable global para asignar conexion
var DB *sql.DB

// Estructura que representa el modelo User
type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Sex       string `json:"sex"`
	Birthdate string `json:"birthdate"`
	Created   string `json:"created"`
}

func main() {
	r := gin.Default()

	DB, err := sql.Open("mysql", "root@/example")
	if err != nil {
		panic(err.Error())
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		panic(err.Error())
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/users", func(c *gin.Context) {
		user := User{}                 // crear estructura donde se guardara el json de usuario
		err := c.ShouldBindJSON(&user) // func que decodifica el body del request en ela estructura y valida que sea un json valido
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		stmt, err := DB.Query("INSERT INTO users (`name`,`sex`,`birthdate`) VALUES (?,?,?)", user.Name, user.Sex, user.Birthdate)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		c.JSON(200, user)
	})

	r.GET("/users", func(c *gin.Context) {
		rows, err := DB.Query("SELECT * FROM users")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var users []User // array donde se guardaran los datos traidos por el query
		for rows.Next() {
			var user User // se crea una var temporal para asignar el valor de la iteracion
			rows.Scan(&user.Id, &user.Name, &user.Sex, &user.Birthdate, &user.Created)
			users = append(users, user) // la variable temporal se 'mete' al array de tados
		}

		c.JSON(200, users)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		var user User // crear estructura donde se guardara el json de usuario
		err := DB.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&user.Id, &user.Name, &user.Sex, &user.Birthdate, &user.Created)
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"msg": "user not found"})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, user)
	})

	r.Run(":8069") // listen and serve on 0.0.0.0:8069 (for windows "localhost:8069")
}
