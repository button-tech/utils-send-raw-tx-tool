package api

import (
	"encoding/json"
	"github.com/button-tech/logger"
	"net/http"
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

func NewServer() (*Server, error) {
	s := server()
	s.initBaseRoute()
	s.fs()
	err := createInfoResponse()
	s.R.Use(cors)

	return s, err
}

func server() *Server {
	return &Server{
		R: routing.New(),
	}
}

func (s *Server) fs() {
	s.Core = &fasthttp.Server{
		Handler:      s.R.HandleRequest,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func cors(ctx *routing.Context) error {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", string(ctx.Request.Header.Peek("Origin")))
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "false")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,POST,DELETE")
	ctx.Response.Header.Set(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	)

	if string(ctx.Method()) == "OPTIONS" {
		ctx.Abort()
	}
	if err := ctx.Next(); err != nil {
		if httpError, ok := err.(routing.HTTPError); ok {
			ctx.Response.SetStatusCode(httpError.StatusCode())
		} else {
			ctx.Response.SetStatusCode(http.StatusInternalServerError)
		}

		b, err := json.Marshal(err)
		if err != nil {
			respondWithJSON(ctx, fasthttp.StatusInternalServerError, map[string]interface{}{
				"error": err},
			)
			logger.Error("cors", err)
		}
		ctx.SetContentType("application/json")
		ctx.SetBody(b)
	}
	return nil
}

func (s *Server) initBaseRoute() {
	s.G = s.R.Group(groupURL)
	s.G.Post("/send", sendHandler)
	s.G.Get("/info", infoHandler)
}

func createInfoResponse() error {
	var err error
	infoResponse, err = json.MarshalIndent(&info, "", "")
	return err
}

func respondWithJSON(ctx *routing.Context, code int, payload interface{}) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(payload); err != nil {
		logger.Error("write answer", err)
	}
}
