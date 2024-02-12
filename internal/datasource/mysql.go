package datasource

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	MySQLDuplicateEntry = 1062
	mysqlDefaultTimeout = time.Second * 30
	mysqlDefaultMaxConn = 30
)

type MySQL struct {
	*sqlx.DB
}

type MySQLConfig struct {
	User    string
	Passwd  string
	Host    string
	Port    string
	DBName  string
	MaxConn int
	Timeout time.Duration
}

func (c *MySQLConfig) FormatDSN() string {
	loc, _ := time.LoadLocation("Asia/Seoul")

	config := &mysql.Config{
		User:                 c.User,
		Passwd:               c.Passwd,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", c.Host, c.Port),
		DBName:               c.DBName,
		Timeout:              c.Timeout,
		ReadTimeout:          c.Timeout,
		WriteTimeout:         c.Timeout,
		Collation:            "utf8mb4_unicode_ci",
		ParseTime:            true,
		Loc:                  loc,
		AllowNativePasswords: true,
		Params: map[string]string{
			"time_zone": fmt.Sprintf("'%s'", loc.String()),
		},
	}
	return config.FormatDSN()
}

func NewMySQL(config *MySQLConfig) *sqlx.DB {
	if config.MaxConn == 0 {
		config.MaxConn = mysqlDefaultMaxConn
	}
	if config.Timeout == 0 {
		config.Timeout = mysqlDefaultTimeout
	}
	db, err := sqlx.Connect("mysql", config.FormatDSN())
	if err != nil {
		panic(err)
	}
	return db
}
