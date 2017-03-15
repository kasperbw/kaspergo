package webserver

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

type webServer struct {
	router     *httprouter.Router
	middleware *negroni.Negroni
	renderer   *render.Render
	port       int
}

var server *webServer

//New webserver instance 생성
func New(listenPort int) {
	server = &webServer{
		router:     httprouter.New(),
		middleware: negroni.Classic(),
		renderer:   render.New(),
		port:       listenPort,
	}

	server.middleware.UseHandler(server.router)
}

//Run webserver run
func Run() {
	if server == nil {
		panic("webserver instance is nil")
	}

	address := fmt.Sprintf(":%d", server.port)
	server.middleware.Run(address)
}

//RegisterRouterHandle router handle 등록
func RegisterRouterHandle(method string, path string, handle httprouter.Handle) {
	switch method {
	case "GET":
		server.router.GET(path, handle)
	case "POST":
		server.router.POST(path, handle)
	default:
		panic(fmt.Sprintf("invalid method: %s", method))
	}
}

//Render 결과 페이지 render
func Render(dataType string, w http.ResponseWriter, status int, name string, value interface{}, htmlOpt ...render.HTMLOptions) {
	switch dataType {
	case "HTML":
		server.renderer.HTML(w, status, name, value, htmlOpt...)
	case "JSON":
		server.renderer.JSON(w, status, value)
	default:
		panic(fmt.Sprintf("invalid data type: %s", dataType))
	}
}

func SetSessionStore(storeType string) {

}
