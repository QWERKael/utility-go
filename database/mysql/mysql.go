package mysql

import (
	"fmt"
	"github.com/QWERKael/utility-go/log"
	"github.com/go-mysql-org/go-mysql/client"
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
	DB *client.Conn
}

func (ci *ConnectInfo) NewConnector() (*Connector, error) {
	conn := &Connector{}
	err := conn.Connect(ci.userName, ci.password, ci.host, ci.port, ci.database)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Connector) Connect(userName string, password string, host string, port int, database string) error {
	conn, err := client.Connect(fmt.Sprintf("%s:%d", host, port), userName, password, database, nil)
	if err != nil {
		log.SugarLogger.Debugf("连接到【%s:%d】失败", host, port)
		return err
	}
	c.DB = conn
	return nil
}

func (c *Connector) QueryAsMapList(sql string) ([]map[string]interface{}, error) {
	r, err := c.DB.Execute(sql)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	colNames := make([]string, r.ColumnNumber())
	for colName, i := range r.FieldNames {
		colNames[i] = colName
	}
	var resultSet = make([]map[string]interface{}, r.RowNumber())
	for i, row := range r.Values {
		resultSet[i] = make(map[string]interface{})
		for j, col := range row {
			resultSet[i][colNames[j]] = col.Value()
		}
	}
	return resultSet, nil
}

func (c *Connector) QueryAsMapStringListWithColNames(sql string) ([]map[string]string, []string, error) {
	r, err := c.DB.Execute(sql)
	if err != nil {
		return nil, nil, err
	}
	defer r.Close()
	colNames := make([]string, r.ColumnNumber())
	for colName, i := range r.FieldNames {
		colNames[i] = colName
	}
	var resultSet = make([]map[string]string, r.RowNumber())
	for i, row := range r.Values {
		resultSet[i] = make(map[string]string)
		for j, col := range row {
			resultSet[i][colNames[j]] = col.String()
		}
	}
	return resultSet, colNames, nil
}
