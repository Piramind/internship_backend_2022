package main

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "bal0nk@"
	dbname   = "postgres"
)

// database immitation
type balancereq struct {
	ID      string `json:"id"`
	Balance int    `json:"balance"`
}

type reservereq struct {
	ID         string `json:"id"`
	Service_id string `json:"service_id"`
	Order_id   string `kson:"order_id"`
	Balance    int    `json:"balance"`
}

func getBalance(c *gin.Context) {
	var balance balancereq //ID, Balance
	id := c.Param("id")
	balance.ID = id
	var b int
	var err error
	b, err = getBalanceById(id)
	if err != nil {
		panic(err)
	}
	balance.Balance = b
	c.IndentedJSON(http.StatusOK, balance)
}

func updateBalance(id string, balance int) {
	sqlStatement := `
	UPDATE users
	SET balance = $2
	WHERE id = $1;`
	_, err := db.Exec(sqlStatement, id, balance)
	if err != nil {
		panic(err)
	}

}

func createBalance(c *gin.Context) {
	var newBal balancereq

	if err := c.BindJSON(&newBal); err != nil {
		return
	}
	sqlStatement := `INSERT INTO users (id, balance) VALUES ($1, $2)`

	_, err := db.Exec(sqlStatement, newBal.ID, newBal.Balance)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, newBal)
}

func createReserve(c *gin.Context) {
	var newRes reservereq

	if err := c.BindJSON(&newRes); err != nil {
		return
	}
	sqlStatement := `
	INSERT INTO reserves (id, service_id, order_id, balance)
	VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(sqlStatement, newRes.ID, newRes.Service_id, newRes.Order_id, newRes.Balance)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, newRes)
}

func addMoney(c *gin.Context) {
	var newReq balancereq

	if err := c.BindJSON(&newReq); err != nil {
		return
	}

	a, _ := getBalanceById(newReq.ID)
	updateBalance(newReq.ID, newReq.Balance+a)
	newReq.Balance = newReq.Balance + a
	c.IndentedJSON(http.StatusOK, newReq)
}

func getBalanceById(id string) (int, error) {

	sqlStatement := `SELECT balance FROM users WHERE id=$1;`
	var money int
	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	row := db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&money); err {
	case sql.ErrNoRows:
		return 0, errors.New("book not found")
	case nil:
		return money, nil
	default:
		panic(err)
	}

}

func getReserveByData(id string, service_id string, order_id string, price int) (int, error) {
	sqlStatement := `SELECT balance FROM reserves WHERE id=$1 and service_id=$2 AND order_id=$3 AND balance=$4;`
	var money int
	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	row := db.QueryRow(sqlStatement, id, service_id, order_id, price)
	switch err := row.Scan(&money); err {
	case sql.ErrNoRows:
		return 0, errors.New("book not found")
	case nil:
		return money, nil
	default:
		panic(err)
	}
}

func updateReserve(id string, service_id string, order_id string, balance int) {
	sqlStatement := `
	UPDATE reserves
	SET balance = $4
	WHERE id = $1 AND service_id=$2 AND order_id=$3;`
	_, err := db.Exec(sqlStatement, id, service_id, order_id, balance)
	if err != nil {
		panic(err)
	}

}

func addReserve(c *gin.Context) {
	var newReq reservereq

	if err := c.BindJSON(&newReq); err != nil {
		return
	}

	a, _ := getReserveByData(newReq.ID, newReq.Service_id, newReq.Order_id, newReq.Balance)
	updateReserve(newReq.ID, newReq.Service_id, newReq.Order_id, newReq.Balance+a)
	newReq.Balance = newReq.Balance + a
	c.IndentedJSON(http.StatusOK, newReq)
}

func admitReserve(c *gin.Context) {
	var newReq reservereq

	if err := c.BindJSON(&newReq); err != nil {
		return
	}

	sqlStatement := `
	DELETE FROM reserves
	WHERE id = $1 AND service_id = $2 AND order_id = $3 AND balance = $4;`
	_, err := db.Exec(sqlStatement, newReq.ID, newReq.Service_id, newReq.Order_id, newReq.Balance)
	if err != nil {
		panic(err)
	}

	a, _ := getBalanceById(newReq.ID)
	updateBalance(newReq.ID, a-newReq.Balance)

	var retReq balancereq
	retReq.ID = newReq.ID
	retReq.Balance = a - newReq.Balance
	c.IndentedJSON(http.StatusOK, retReq)
}

func connectDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db
}

var db = connectDatabase()

func main() {
	defer db.Close()

	router := gin.Default()
	router.GET("/balances/:id", getBalance)       //yes
	router.POST("/create_balance", createBalance) //yes
	router.POST("/create_reserve", createReserve) //yes
	router.POST("/transfer", addMoney)            //yes
	router.POST("/reserve", addReserve)           //yes
	router.POST("/admit", admitReserve)
	router.Run("localhost:8080") //yes
}
