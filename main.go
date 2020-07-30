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
	} else {
		logrus.Debugf("From DeleteOneUsersecAllAccess :  %v", result)
	}
}

// CreateOneUsersec creates one usersec
func CreateOneUsersec(c *gin.Context) {

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	dbpool := createConnection()
	defer dbpool.Close()

	var usersec Usersec
	err := c.ShouldBind(&usersec)
	if err != nil {
		panic(err)
	}

	db := createConnection()
	defer db.Close()

	myQuery := `INSERT INTO usersec (userid, menuname, mainmenu) VALUES ($1, $2, $3)`

	result, err := dbpool.Exec(c, myQuery, usersec.UserID, usersec.Menuname, usersec.Mainmenu)

	if err != nil {
		panic(err)
	} else {
		logrus.Debugf("From CreateOneUsersec :  %v", result)
	}
}






// Login struct (Model)
type Login struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password1 string `json:"password"`
}

// Init allLogins var as a slice Login struct
var allLogins []Login

// GetLogins returns allLogins
func GetLogins(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	dbpool := createConnection()
	defer dbpool.Close()

	myQuery := `SELECT * FROM logins`

	rows, err := dbpool.Query(c, myQuery)
	if err != nil {
		panic(err)
	} else {
		defer rows.Close()

		var username string
		var email string
		var password1 string

		allLogins = nil
		for rows.Next() {
			rows.Scan(&username, &email, &password1)
			allLogins = append(allLogins, Login{Username: username, Email: email, Password1: password1})
		}

		if rows.Err() != nil {
			panic(rows.Err())
		} else {
			c.JSON(200, allLogins)
		}
	}
}

// GetLoginUsername returns single login by username
func GetLoginUsername(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	dbpool := createConnection()
	defer dbpool.Close()

	usernameFromRouter := c.Param("username")

	myQuery := `SELECT * FROM logins WHERE username=$1`

	rows, err := dbpool.Query(c, myQuery, usernameFromRouter)
	if err != nil {
		panic(err)
	} else {
		defer rows.Close()

		var username string
		var email string
		var password1 string

		allLogins = nil
		for rows.Next() {
			rows.Scan(&username, &email, &password1)
			allLogins = append(allLogins, Login{Username: username, Email: email, Password1: password1})
		}

		if rows.Err() != nil {
			panic(rows.Err())
		} else {
			c.JSON(200, allLogins)
		}
	}

}

// GetLoginEmail returns single login by email
func GetLoginEmail(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	dbpool := createConnection()
	defer dbpool.Close()

	emailFromRouter := c.Param("email")

	myQuery := `SELECT * FROM logins WHERE email=$1`

	rows, err := dbpool.Query(c, myQuery, emailFromRouter)
	if err != nil {
		panic(err)
	} else {
		defer rows.Close()

		var username string
		var email string
		var password1 string

		allLogins = nil
		for rows.Next() {
			rows.Scan(&username, &email, &password1)
			allLogins = append(allLogins, Login{Username: username, Email: email, Password1: password1})
		}

		if rows.Err() != nil {
			panic(rows.Err())
		} else {
			c.JSON(200, allLogins)
		}
	}
}

// CreateLogin Adds new login
func CreateLogin(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	dbpool := createConnection()
	defer dbpool.Close()

	var login Login

	err := c.ShouldBind(&login)
	if err != nil {
		panic(err)
	} else {

		db := createConnection()
		defer db.Close()

		myQuery := `INSERT INTO logins (username, email, password1) VALUES ($1, $2, $3)`

		result, err := dbpool.Exec(c, myQuery, login.Username, login.Email, login.Password1)

		if err != nil {
			panic(err)
		} else {
			logrus.Debugf("From CreateLogin :  %v", result)
		}
	}
}






// Purchase struct (Model)
type Purchase struct {
	ID       string `json:"item_id"`
	Name     string `json:"item_name"`
	Quantity string `json:"item_quantity"`
	Rate     string `json:"item_rate"`
	Date     string `json:"item_purchase_date"`
}

// Init allpurchases var as a slice Purchase struct
var allPurchases []Purchase

// GetPurchases returns all purchases
func GetPurchases(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	dbpool := createConnection()
	defer dbpool.Close()

	myQuery := `SELECT * FROM purchases`

	rows, err := dbpool.Query(c, myQuery)
	if err != nil {
		panic(err)
	} else {
		defer rows.Close()

		var id string
		var name string
		var quantity string
		var rate string
		var date string

		allPurchases = nil
		for rows.Next() {
			rows.Scan(&id, &name, &quantity, &rate, &date)
			allPurchases = append(allPurchases, Purchase{ID: id, Name: name, Quantity: quantity, Rate: rate, Date: date})
		}

		if rows.Err() != nil {
			panic(rows.Err())
		} else {
			c.JSON(200, allPurchases)
		}
	}
}

// GetPurchase returns one purchase
func GetPurchase(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	dbpool := createConnection()
	defer dbpool.Close()

	idFromRouter := c.Param("id")

	myQuery := `SELECT * FROM purchases WHERE item_id=$1`

	rows, err := dbpool.Query(c, myQuery, idFromRouter)
	if err != nil {
		panic(err)
	} else {
		defer rows.Close()

		var id string
		var name string
		var quantity string
		var rate string
		var date string

		allPurchases = nil
		for rows.Next() {
			rows.Scan(&id, &name, &quantity, &rate, &date)
			allPurchases = append(allPurchases, Purchase{ID: id, Name: name, Quantity: quantity, Rate: rate, Date: date})
		}

		if rows.Err() != nil {
			panic(rows.Err())
		} else {
			c.JSON(200, allPurchases)
		}
	}
}

// DeletePurchase deletes one purchase
func DeletePurchase(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	dbpool := createConnection()
	defer dbpool.Close()

	idFromRouter := c.Param("id")

	myQuery := `DELETE FROM purchases WHERE item_id=$1`

	result, err := dbpool.Exec(c, myQuery, idFromRouter)
	if err != nil {
		panic(err)
	} else {
		logrus.Debugf("From DeletePurchase :  %v", result)
	}
}

// CreatePurchase adds new purchase
func CreatePurchase(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	dbpool := createConnection()
	defer dbpool.Close()

	var puchase Purchase

	err := c.ShouldBind(&puchase)
	if err != nil {
		panic(err)
	} else {

		db := createConnection()
		defer db.Close()

		myQuery := `INSERT INTO purchases (item_id, item_name, item_quantity, item_rate, item_purchase_date) VALUES ($1, $2, $3, $4, $5)`

		result, err := dbpool.Exec(c, myQuery, puchase.ID, puchase.Name, puchase.Quantity, puchase.Rate, puchase.Date)

		if err != nil {
			panic(err)
		} else {
			logrus.Debugf("From CreatePurchase :  %v", result)
		}
	}
}

// UpdatePurchase updates new purchase
func UpdatePurchase(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	dbpool := createConnection()
	defer dbpool.Close()

	var puchase Purchase

	err := c.ShouldBind(&puchase)
	if err != nil {
		panic(err)
	} else {

		db := createConnection()
		defer db.Close()

		idFromRouter := c.Param("id")

		myQuery := `UPDATE purchases SET item_id=$1, item_name=$2, item_quantity=$3, item_rate=$4, item_purchase_date=$5 WHERE item_id=$6`

		result, err := dbpool.Exec(c, myQuery, puchase.ID, puchase.Name, puchase.Quantity, puchase.Rate, puchase.Date, idFromRouter)

		if err != nil {
			panic(err)
		} else {
			logrus.Debugf("From UpdatePurchase :  %v", result)
		}
	}
}




func main() {
	router := gin.Default()



	router.GET("/usersec", GetAllUsersec)
	router.GET("/usersec/:userid", GetOneUsersec)
	router.DELETE("/usersec/:userid", DeleteOneUsersecAllAccess)
	router.POST("/usersec", CreateOneUsersec)






	router.GET("/logins", GetLogins)
	router.GET("/logins/username/:username", GetLoginUsername)
	router.GET("/logins/email/:email", GetLoginEmail)
	router.POST("/logins", CreateLogin)





	router.GET("/purchases", GetPurchases)
	router.GET("/purchases/:id", GetPurchase)
	router.DELETE("/purchases/:id", DeletePurchase)
	router.PUT("/purchases/:id", UpdatePurchase)
	router.POST("/purchases", CreatePurchase)






	router.Use(cors.Default())
	router.Run(":" + goDotEnvVariable("API_PORT"))
}
