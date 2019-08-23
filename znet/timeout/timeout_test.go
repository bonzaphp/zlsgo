package timeout

import (
	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/znet"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWebTimeout(T *testing.T) {
	t := zlsgo.NewTest(T)
	r := newServer()

	w1 := newRequest(r, "GET", "/timeout_1", New(1*time.Second), func(c *znet.Context) {
		c.String(200, "timeout_1")
	})
	t.Equal(200, w1.Code)
	t.Equal("timeout_1", w1.Body.String())

	w2 := newRequest(r, "GET", "/timeout_2", New(1*time.Second), func(c *znet.Context) {
		time.Sleep(2 * time.Second)
		c.String(200, "timeout_2")
	})
	t.Equal(504, w2.Code)
	t.Equal("", w2.Body.String())

	w3 := newRequest(r, "GET", "/timeout_3", New(1*time.Second, func(c *znet.Context) {
		c.String(210, "is timeout")
	}), func(c *znet.Context) {
		time.Sleep(2 * time.Second)
		c.String(200, "timeout_3")
	})
	t.Equal(210, w3.Code)
	t.Equal("is timeout", w3.Body.String())
}

var (
	one    sync.Once
	Engine *znet.Engine
)

func newServer() *znet.Engine {
	one.Do(func() {
		Engine = znet.New()
		Engine.SetMode(znet.DebugMode)
	})
	return Engine
}

func newRequest(r *znet.Engine, method string, path string, handler ...znet.HandlerFunc) *httptest.ResponseRecorder {
	method = strings.ToUpper(method)
	if len(handler) > 0 {
		firstHandler := handler[0]
		handlers := handler[1:]
		switch method {
		case "GET":
			r.GET(path, firstHandler, handlers...)
		case "POST":
			r.POST(path, firstHandler, handlers...)
		}
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}