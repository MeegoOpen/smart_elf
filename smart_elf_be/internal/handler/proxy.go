package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"smart_elf_standalone/internal/auth"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	target     *url.URL
	feishuAuth *auth.FeishuAuth
}

func NewProxyHandler(target string, feishuAuth *auth.FeishuAuth) *ProxyHandler {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	return &ProxyHandler{
		target:     targetURL,
		feishuAuth: feishuAuth,
	}
}

func (h *ProxyHandler) ProxyRequest(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(h.target)
	proxy.Director = func(req *http.Request) {
		req.Host = h.target.Host
		req.URL.Scheme = "https"
		req.URL.Host = h.target.Host
		req.URL.Path = c.Param("path")

		token, err := h.feishuAuth.GetToken()
		if err != nil {
			// Handle error, maybe return an error to the client
			return
		}
		req.Header.Set("X-Plugin-Token", token)
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
