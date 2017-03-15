package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"core/mssql"
	"core/webserver"
)

func main() {
	mssql.Initialize("./config/dbconfig.json")

	chan1 := make(chan int, 1)
	chan2 := make(chan int, 1)

	go queryTest(2, 5013, chan1)
	go queryTest(1, 5002, chan2)

	<-chan1
	<-chan2
}

func webserverTest() {
	webserver.New(3000)

	registerRouter()

	webserver.Run()
}

func registerRouter() {
	webserver.RegisterRouterHandle("GET", "/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		webserver.Render("HTML", w, http.StatusOK, "index", map[string]string{"title": "Test Server"})
	})
}

func queryTest(index int, userSeq int, reportCh chan int) {
	for i := 0; i < 1000; i++ {
		result, err := mssql.Query("userdb", index, "execute query", userSeq)
		if err != nil {
			panic(err)
		}

		fmt.Printf("count(%d) : %v\n", i, result.Data[0][0]["user_seq"].(int64))
	}

	reportCh <- 1
	/*
		var b bytes.Buffer

		if err := json.NewEncoder(&b).Encode(result.Data); err != nil {
			panic(err)
		}

		fo, err := os.Create("e:\\iw2admin\\src\\spLoadTable_result3.json")
		if err != nil {
			panic(err)
		}
		defer fo.Close()

		buf := b.Bytes()

		if _, err := fo.Write(buf); err != nil {
			panic(err)
		}
	*/
}
