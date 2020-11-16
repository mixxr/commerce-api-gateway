package router

import "net/http"

type IRouter interface {
	GetHandler() http.Handler
}
