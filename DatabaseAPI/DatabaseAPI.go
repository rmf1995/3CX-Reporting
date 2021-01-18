package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

// C config var from strut
var C config

type config struct {
	Database database
}

type database struct {
	User string `mapstructure:"User"`
	Pass string `mapstructure:"Pass"`
	Name string `mapstructure:"Name"`
	Host string `mapstructure:"Host"`
	Port string `mapstructure:"Port"`
}

type server struct {
	ID         string `json:"id"`
	Name       string `json:"Name"`
	Location   string `json: "Location"`
	URL        string `json : "URL"`
	UserName   string `json : "UserName"`
	PwStateID  string `json : "PwStateID"`
	CustomerID string `json : "CustomerID"`
}
type servers struct {
	servers []server `json:"employee"`
}

func dbConn() (db *sql.DB) {

	dbDriver := "mysql"
	dbUser := C.Database.User
	dbPass := C.Database.Pass
	dbName := C.Database.Name
	dbHost := C.Database.Host
	dbPort := C.Database.Port

	connStr := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
	db, err := sql.Open(dbDriver, connStr)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// Hello default funcion when / is called.
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Access Denied.")
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print("[" + time.Now().UTC().Format(time.RFC3339) + "] " + string(bytes))
}

func getAll3CX(c echo.Context) error {
	db := dbConn()
	rows, err := db.Query("SELECT * FROM servers")
	if err != nil {
		//
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		//
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		//
	}
	defer db.Close()
	return c.String(http.StatusOK, string(jsonData))
}

func delete3CX(c echo.Context) error {
	db := dbConn()

	requestedid := c.Param("id")
	sql := "DELETE FROM servers WHERE id = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println(err)
	}
	_, err2 := stmt.Exec(requestedid)
	if err2 != nil {
		log.Println(err2)
	}
	defer db.Close()
	return c.JSON(http.StatusOK, "Deleted")
}

func create3CX(c echo.Context) error {
	db := dbConn()

	server := new(server)
	if err := c.Bind(server); err != nil {
		return err
	}
	//
	sql := "INSERT INTO servers (Name, Location, URL, UserName, PwStateID, CustomerID) VALUES ( ?, ?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(sql)

	if err != nil {
		fmt.Print(err.Error())
	}
	defer stmt.Close()
	result, err2 := stmt.Exec(server.Name, server.Location, server.URL, server.UserName, server.PwStateID, server.CustomerID)

	// Exit if we get an error
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(result.LastInsertId())
	defer db.Close()
	return c.JSON(http.StatusCreated, server.Name)
}

func edit3CX(c echo.Context) error {
	db := dbConn()

	serverID := c.Param("id")
	var id string
	var name string
	var location string
	var url string
	var username string
	var pwstateID string
	var customerID string

	err := db.QueryRow("SELECT id,Name,Location,URL,UserName,PwStateID,CustomerID FROM servers WHERE id = ?", serverID).Scan(&id, &name, &location, &url, &username, &pwstateID, &customerID)
	if err != nil {
		//log.Println(err)
	}

	response := server{ID: id, Name: name, Location: location, URL: url, UserName: username, PwStateID: pwstateID, CustomerID: customerID}
	defer db.Close()
	return c.JSON(http.StatusOK, response)
}

func update3CX(c echo.Context) error {
	db := dbConn()

	requestedid := c.Param("id")
	server := new(server)
	if err := c.Bind(server); err != nil {
		return err
	}
	log.Println(server)

	sql := "UPDATE servers SET Name=?, Location=?, URL=?, UserName=?, PwStateID=?, CustomerID=? WHERE id=?"
	stmt, err := db.Prepare(sql)

	if err != nil {
		log.Println(err.Error())
	}

	defer stmt.Close()

	result, err2 := stmt.Exec(server.Name, server.Location, server.URL, server.UserName, server.PwStateID, server.CustomerID, requestedid)
	log.Println(result)
	// Exit if we get an error
	if err2 != nil {
		log.Println(err2)
	}

	defer db.Close()
	return c.JSON(http.StatusCreated, server.Name)
}

func main() {
	log.Printf("3CX Reporting Tool (DatabaseAPI) by rmf1995")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/opt/3CX-Reporting")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("Unable to read config file")
		viper.GetViper().ConfigFileUsed()

		os.Exit(1)
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("Unable to Decode into Struct, %v", err)
	}

	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	e := echo.New()
	e.HideBanner = true
	e.Debug = false
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.GET("/", Hello)
	e.GET("/v1/getAll3CX/", getAll3CX)
	e.DELETE("/v1/3cx/delete/:id", delete3CX)
	e.POST("/v1/3cx/create", create3CX)
	e.GET("/v1/3cx/edit/:id", edit3CX)
	e.POST("/v1/3cx/update/:id", update3CX)

	s := &http.Server{
		Addr:         ":3055",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	e.StartServer(s)

}
