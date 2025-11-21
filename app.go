package mini

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type App struct{}
type Request = *http.Request
type Response = http.ResponseWriter

var router = http.NewServeMux()
var middlewareList []func(http.Handler) func(Response, Request)

type Contex struct {
	Request
	Response
}

func NewApp() *App {
	return &App{}
}
func InitCtx() *Contex {
	return &Contex{}
}

func (a *App) Listen(port string, callback string) {
	fmt.Println(callback)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) Use(middleware func(http.Handler) func(Response, Request)) {
	middlewareList = append(middlewareList, middleware)
}
func chainMiddleware(handler func(Response, Request)) func(Response, Request) {
	switch len(middlewareList) {
	case 1:
		return middlewareList[0](http.HandlerFunc(handler))
	case 2:
		return middlewareList[0](http.HandlerFunc(
			middlewareList[1](http.HandlerFunc(handler))))

	case 3:
		return middlewareList[0](http.HandlerFunc(
			middlewareList[1](http.HandlerFunc(
				middlewareList[2](http.HandlerFunc(handler))))))

	case 4:
		return middlewareList[0](http.HandlerFunc(
			middlewareList[1](http.HandlerFunc(
				middlewareList[2](http.HandlerFunc(
					middlewareList[3](http.HandlerFunc(handler))))))))

	default:
		return handler
	}
}

func (a *App) Get(path string, handler func(Response, Request)) {
	if len(middlewareList) != 0 {
		router.HandleFunc("GET "+path, chainMiddleware(handler))
	} else {
		router.HandleFunc("GET "+path, handler)
	}
}

func (a *App) Post(path string, handler func(Response, Request)) {
	router.HandleFunc("POST "+path, chainMiddleware(handler))
}

func (a *App) Put(path string, handler func(Response, Request)) {
	router.HandleFunc("PUT "+path, chainMiddleware(handler))
}
func (a *App) Delete(path string, handler func(Response, Request)) {
	router.HandleFunc("DELETE "+path, chainMiddleware(handler))
}

// Context Methods
func (c *Contex) Add(res Response, req Request) {
	c.Request = req
	c.Response = res
}

func (c *Contex) Send(msg string) {
	fmt.Fprint(c.Response, msg)
}

func (c *Contex) Json(body map[string]string) {
	json.NewEncoder(c.Response).Encode(body)
}

func (c *Contex) Body(value interface{}) {
	err := json.NewDecoder(c.Request.Body).Decode(value)
	if err != nil {
		http.Error(c.Response, "Error Decoding Json", http.StatusBadRequest)
	}

	defer c.Request.Body.Close()
}

func (c *Contex) Params(value string) string {
	return c.Request.PathValue(value)
}

// func (c *Contex) myParams() string {
// 	urlPath := c.Request.URL.Path

// 	newUrl, err := url.Parse(urlPath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return path.Base(newUrl.Path)

// }
