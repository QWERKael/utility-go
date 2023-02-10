package mysql

import (
	"database/sql"
	"fmt"
	"github.com/QWERKael/utility-go/log"
	_ "github.com/go-sql-driver/mysql"
)

type ConnectInfo struct {
	userName string
	password string
	network  string
	host     string
	port     int
	database string
}

type Connector struct {
	DB *sql.DB
}

func NewConnector(ci *ConnectInfo) (*Connector, error) {
	conn := &Connector{}
	err := conn.Connect(ci.userName, ci.password, ci.network, ci.host, ci.port, ci.database)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Connector) Connect(userName string, password string, network string, host string, port int, database string) error {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", userName, password, network, host, port, database)
	return c.ConnectWithDSN(dsn)
}

func (c *Connector) ConnectWithDSN(dsn string) error {
	log.SugarLogger.Debugf("连接到【%s】", dsn)
	var err error
	c.DB, err = sql.Open("mysql", dsn)
	return err
}

func (c *Connector) QueryAsMapList(sql string) ([]map[string]interface{}, error) {
	rows, err := c.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var colNames []string
	colNames, err = rows.Columns()
	if err != nil {
		return nil, err
	}
	var resultSet = make([]map[string]interface{}, 0)

	for rows.Next() {
		cols := make([]interface{}, len(colNames))
		colPtrs := make([]interface{}, len(colNames))
		for i, _ := range cols {
			colPtrs[i] = &cols[i]
		}
		if err := rows.Scan(colPtrs...); err != nil {
			return nil, err
		}
		result := make(map[string]interface{})
		for i, colName := range colNames {
			v := colPtrs[i].(*interface{})
			result[colName] = *v
		}
		resultSet = append(resultSet, result)
	}
	return resultSet, nil
}

func (c *Connector) QueryAsMapStringListWithColNames(sql string) ([]map[string]string, []string, error) {
	rows, err := c.DB.Query(sql)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var colNames []string
	colNames, err = rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	var resultSet = make([]map[string]string, 0)

	for rows.Next() {
		cols := make([]interface{}, len(colNames))
		colPtrs := make([]interface{}, len(colNames))
		for i, _ := range cols {
			colPtrs[i] = &cols[i]
		}
		if err := rows.Scan(colPtrs...); err != nil {
			return nil, nil, err
		}
		result := make(map[string]string)
		for i, colName := range colNames {
			v := colPtrs[i].(*interface{})
			switch {
			case *v == nil:
				result[colName] = "NULL"
				continue
			default:
				result[colName] = fmt.Sprintf("%s", *v)
			}
		}
		resultSet = append(resultSet, result)
	}
	return resultSet, colNames, nil
}
