package mssql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb" //mssql을 사용하기 위해 import해야함.
)

//DbConfigManager singleton 객체로 DbConfigManagerInstance() 호출하여 사용
type DbConfigManager struct {
	Databases []DbConfig `json:"config"`
}

//DbConfig db 설정
type DbConfig struct {
	DbType            string `json:"db_type"`
	Index             int    `json:"db_index"`
	Hostname          string `json:"hostname"`
	Port              int    `json:"port"`
	Database          string `json:"database"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	ConnectionTimeOut int    `json:"timeout"`
}

var dbConnectionManager map[string]map[int]*sql.DB

//Initialize mssql을 사용하기 위해 초기화
func Initialize(configfile string) {
	file, err := ioutil.ReadFile(configfile)
	if err != nil {
		panic(err)
	}

	var configMgr DbConfigManager
	if err := json.Unmarshal(file, &configMgr); err != nil {
		panic(err)
	}

	dbConnectionManager = map[string]map[int]*sql.DB{}

	for _, v := range configMgr.Databases {
		cs := v.connectionString()
		conn, err := sql.Open("mssql", cs)
		if err != nil {
			panic(err)
		}

		if err := conn.Ping(); err != nil {
			panic(err)
		}

		if _, exists := dbConnectionManager[v.DbType]; !exists {
			dbConnectionManager[v.DbType] = map[int]*sql.DB{}
		}

		dbConnectionManager[v.DbType][v.Index] = conn
	}
}

func connection(dbType string, index int) (*sql.DB, error) {
	if _, exists := dbConnectionManager[dbType]; !exists {
		return nil, fmt.Errorf("Not Found DB : %s", dbType)
	}

	if _, exists := dbConnectionManager[dbType][index]; !exists {
		return nil, fmt.Errorf("Not Found DB : %s(%d)", dbType, index)
	}

	return dbConnectionManager[dbType][index], nil
}

/*
//GetConnectionString connection string 생성
func getConnectionString(dbType string, index int) string {
	for _, v := range configMgr.Databases {
		if v.DbType == dbType && v.Index == index {
			return v.connectionString()
		}
	}

	return ""
}
*/

//ConnectionString db연결 스트링 생성
func (config *DbConfig) connectionString() string {
	query := url.Values{}
	query.Add("database", config.Database)
	query.Add("connection timeout", fmt.Sprintf("%d", config.ConnectionTimeOut))

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(config.Username, config.Password),
		Host:     fmt.Sprintf("%s:%d", config.Hostname, config.Port),
		RawQuery: query.Encode(),
	}

	return u.String()
}
