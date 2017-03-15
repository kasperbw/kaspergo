package mssql

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb" //mssql을 사용하기 위해 import해야함.
)

//QueryResult 쿼리 결과값
type QueryResult struct {
	Count int // resultset 갯수
	Data  map[int][]map[string]interface{}
}

//rowBind sql에서 row결과 값을 받아오기 위한 구조체
type rowBind struct {
	columns  []string
	datas    []interface{} //실제 결과 값을 받을 배열
	dataPtrs []interface{} //결과 값 받을 배열 각각의 포인트

	result []map[string]interface{}
}

func (rb *rowBind) prepare(rows *sql.Rows) error {
	var err error
	//row의 column 이름 배열 구하기
	rb.columns, err = rows.Columns()
	if err != nil {
		return err
	}

	//row에서 반환 하는 column 수 만큼 저장을 위한 배열 생성
	rb.datas = make([]interface{}, len(rb.columns))
	rb.dataPtrs = make([]interface{}, len(rb.columns))

	//실제 data가 저장될 포인트를 저장
	for index := range rb.dataPtrs {
		rb.dataPtrs[index] = &rb.datas[index]
	}

	return nil
}

func (rb *rowBind) makeResult() {
	//현재 row의 data값이 column과 매칭되어 저장될 map
	m := make(map[string]interface{}, len(rb.columns))

	for index, column := range rb.columns {
		m[column] = rb.datas[index]
	}

	rb.result = append(rb.result, m)
}

//Exec 결과값이 없을 경우 사용. 적용된 row갯수를 반환한다.
func Exec(dbName string, dbIndex int, query string, params ...interface{}) (int64, error) {
	/*
		connString := getConnectionString(dbName, dbIndex)

		//sql connection 열기
		conn, err := sql.Open("mssql", connString)
		if err != nil {
			return 0, err
		}
		defer conn.Close()
	*/
	conn, err := connection(dbName, dbIndex)
	if err != nil {
		return 0, err
	}

	stmt, err := conn.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(params...)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return n, nil
}

//Query 결과 값을 받을 경우 사용
func Query(dbName string, dbIndex int, query string, params ...interface{}) (*QueryResult, error) {
	/*
		connString := getConnectionString(dbName, dbIndex)

		//sql connection 열기
		conn, err := sql.Open("mssql", connString)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
	*/

	conn, err := connection(dbName, dbIndex)
	if err != nil {
		return nil, err
	}

	stmt, err := conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//최종 결과값이 저장될 객체 생성
	result := &QueryResult{
		Count: 0,
		Data:  map[int][]map[string]interface{}{},
	}

	//다중 select처리를 위한 무한 루프
	for {
		//해당 row의 결과값 저장을 위한 객체 준비
		var rb = rowBind{}
		if err := rb.prepare(rows); err != nil {
			return nil, err
		}

		//select 값 받기
		for rows.Next() {
			//error가 sql.ErrNoRows 이면 결과값이 없는 것이다. 에러는 아님.
			if err := rows.Scan(rb.dataPtrs...); err != nil && err != sql.ErrNoRows {
				return nil, err
			}

			rb.makeResult()
		}

		result.Data[result.Count] = rb.result
		result.Count++

		//다음 select 결과 set이 있는지 검사하고 없으면 최종 결과 리턴
		if rows.NextResultSet() == false {
			break
		}
	}

	return result, nil
}
