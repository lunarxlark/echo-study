package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var db *gorm.DB
var err error

type user struct {
	User_id      int    `json:"user_id"`
	Name_first   string `json:"name_first"`
	Name_familly string `json:"name_familly"`
}

func main() {

	dbInit()
	e := echo.New()

	// log
	//file, err := os.OpenFile("access.log", os.O_WRONLY|os.O_CREATE, 0644)
	file, err := os.OpenFile("./log/access.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	// Routing
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Path Parameters
	e.GET("/user/:id", getUser)
	// Query Parameters
	e.GET("/show", show)

	// default
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           ltsvLogFormat(),
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           file,
	}))

	e.Logger.Fatal(e.Start(":1323"))
}

func dbInit() {
	connect := "eapp:password@tcp(db:3306)/esample"
	db, err = gorm.Open("mysql", connect)
	// defer db.Close()
	if err != nil {
		panic(err.Error())
	}
}

func dbConnection() *gorm.DB {
	return db
}

func getUser(c echo.Context) error {
	d := dbConnection()
	id := c.Param("id")
	u := user{}
	u.User_id, _ = strconv.Atoi(id)
	d.First(&u)
	return c.JSON(http.StatusOK, u)
}

func show(c echo.Context) error {
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

func ltsvLogFormat() string {
	var format string
	//format += "time:${time_rfc3339}\t"
	//format += "time:${time_unix}\t"
	format += "time:${time_custom}\t"
	//TODO : LoggerConfigのidが何を表すか
	//format += "id:${id}\t"
	format += "host:${host}\t"
	format += "remote_ip:${remote_ip}\t"
	format += "uri:${uri}\t"
	format += "method:${method}\t"
	//reqのURL.path
	format += "path:${path}\t"
	format += "protocol:${protocol}\t"
	format += "referer:${referer}\t"
	format += "UA:${user-agent}\t"
	format += "status:${status}\t"
	format += "error:${error}\t"
	//format += "latency:${latency}\t"
	format += "latency:${latency}\t"
	format += "latency_human:${latency_human}\t"
	format += "byte_in:${byte_in}\t"
	format += "byte_out:${byte_out}\n"
	return format
}

func tsvLogFormat() string {
	var format string
	format += "${time_rfc3339}\t"
	format += "${host}\t"
	format += "${remote_ip}\t"
	format += "${uri}\t"
	format += "${method}\t"
	format += "${path}\t"
	format += "${protocol}\t"
	format += "${referer}\t"
	format += "${user-agent}\t"
	format += "${status}\t"
	//format += "${error}\t"
	format += "${latency}\t"
	format += "${latency_human}\t"
	format += "${byte_in}\t"
	format += "${byte_out}\n"
	return format
}