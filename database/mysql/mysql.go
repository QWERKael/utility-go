package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Connector struct {
	DB *sql.DB
}

func (c *Connector)Connect(userName string, password string, network string, host string, port int, database string) error {
	connectString := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", userName, password, network, host, port, database)
	var err error
	c.DB, err = sql.Open("mysql", connectString)
	return err
}

func (c *Connector)Query(query string) error {
	rows, err := c.DB.Query(query)
	if err != nil {
		return err
	}
}