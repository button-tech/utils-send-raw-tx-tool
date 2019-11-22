package api

import (
	"encoding/json"
	"log"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const (
	readTimeout  = time.Second * 30
	writeTimeout = time.Second * 30

	groupURL = "/api/v1"
)

type Server struct {
	Core *fasthttp.Server
	R    *routing.Router
	G    *routing.RouteGroup
}

func NewServer() *Server {
	s := server()
	s.initBaseRoute()
	s.fs()
	return s
}

func server() *Server {
	return &Server{
		r: routing.New(),
	}
}

func (s *Server) fs() {
	s.Core = &fasthttp.Server{
		Handler:      s.R.HandleRequest,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func (s *Server) initBaseRoute() {
	s.G = s.R.Group(groupURL)
	s.G.Post("/send", sendHandler)
}

func respondWithJSON(ctx *routing.Context, code int, payload interface{}) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(payload); err != nil {
		log.Println("write answer", err)
	}
}
