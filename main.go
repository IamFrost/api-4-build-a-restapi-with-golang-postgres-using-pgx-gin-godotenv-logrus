package main

import (
	"context"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var logrusLog = logrus.New()

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		logrusLog.Errorf("Error loading .env file %v\n", err)
	}

	return os.Getenv(key)
}

// create database connection
func createConnection() *pgxpool.Pool {
	logrusLog.Infof("Connecting....")
	dbpool, err := pgxpool.Connect(context.Background(), goDotEnvVariable("DATABASE_URL"))
	if err != nil {
		logrusLog.Errorf("Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

// Usersec struct (Model)
type Usersec struct {
	UserID   string `json:"userid"`
	Menuname string `json:"menuname"`
	Mainmenu string `json:"mainmenu"`
}

// Init allUserSec var as a slice Purchase struct
var allUserSec []Usersec

// JSONMiddleware1 returns all header
func JSONMiddleware1() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
		c.Writer.Header().Set("Content-Type", "application/json")
	}
}

// JSONMiddleware2 returns all header
func JSONMiddleware2() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}
}

// GetAllUsersec returns all usersec
func GetAllUsersec(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	dbpool := createConnection()
	defer dbpool.Close()

	myQuery := `SELECT * FROM usersec`

	rows, err := dbpool.Query(c, myQuery)
	if err != nil {
		// c.AbortWithStatus(404)
		panic(err)
		// logrusLog.Errorf("Unable to execute query %v   :   %v", myQuery, err)
		// os.Exit(1)
	} else {
		defer rows.Close()

		var userid string
		var menuname string
		var mainmenu string

		allUserSec = nil
		for rows.Next() {
			rows.Scan(&userid, &menuname, &mainmenu)
			allUserSec = append(allUserSec, Usersec{UserID: userid, Menuname: menuname, Mainmenu: mainmenu})
		}

		if rows.Err() != nil {
			panic(rows.Err())
		} else {
			c.JSON(200, allUserSec)
		}
	}
}

// GetOneUsersec returns one usersec
func GetOneUsersec(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	dbpool := createConnection()
	defer dbpool.Close()

	userid := c.Param("userid")

	myQuery := `SELECT * FROM usersec where userid = $1`

	rows, err := dbpool.Query(c, myQuery, userid)
	if err != nil {
		panic(err)
	} else {
		defer rows.Close()

		var userid string
		var menuname string
		var mainmenu string

		allUserSec = nil
		for rows.Next() {
			rows.Scan(&userid, &menuname, &mainmenu)
			allUserSec = append(allUserSec, Usersec{UserID: userid, Menuname: menuname, Mainmenu: mainmenu})
		}
		// Any errors encountered by rows.Next or rows.Scan will be returned here
		if rows.Err() != nil {
			panic(rows.Err())
		} else {
			c.JSON(200, allUserSec)
		}
	}
}

// DeleteOneUsersecAllAccess deletes one usersec - all access
func DeleteOneUsersecAllAccess(c *gin.Context) {

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	dbpool := createConnection()
	defer dbpool.Close()

	userid := c.Param("userid")

	myQuery := `DELETE FROM usersec WHERE userid = $1`

	result, err := dbpool.Exec(c, myQuery, userid)

	if err != nil {
		panic(err)
	}else {
		logrus.Debugf("From DeleteOneUsersecAllAccess :  %v",result)
	}
}

// CreateOneUsersec creates one usersec
func CreateOneUsersec(c *gin.Context) {

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	dbpool := createConnection()
	defer dbpool.Close()

	var usersec Usersec
	err := c.ShouldBind(&usersec)

	db := createConnection()
	defer db.Close()

	myQuery := `INSERT INTO usersec (userid, menuname, mainmenu) VALUES ($1, $2, $3)`

	result, err := dbpool.Exec(c, myQuery, usersec.UserID, usersec.Menuname, usersec.Mainmenu)

	if err != nil {
		panic(err)
	}else {
		logrus.Debugf("From CreateOneUsersec :  %v",result)
	}
}

func main() {
	router := gin.Default()
	router.GET("/usersec", GetAllUsersec)
	router.GET("/usersec/:userid", GetOneUsersec)
	router.DELETE("/usersec/:userid", DeleteOneUsersecAllAccess)
	router.POST("/usersec", CreateOneUsersec)
	router.Use(cors.Default())
	router.Run(":" + goDotEnvVariable("API_PORT"))
	// router.Use(JSONMiddleware())
}
