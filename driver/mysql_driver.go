package driver

import (
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/spf13/viper"
)

// MysqlDB struct
type MysqlDB struct {
	SQL *gorm.DB
}

// MysqlConnection avairable
var MysqlConnection = &MysqlDB{}

// Connect func
func Connect() (*MysqlDB, error) {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	strConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", strConnection, val.Encode())
	db, err := gorm.Open(`mysql`, dsn)
	MysqlConnection.SQL = db
	return MysqlConnection, err
}
