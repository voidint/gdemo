package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type options struct {
	Port      int
	TargetURL string
}

var opts options

func init() {
	flag.StringVar(&opts.TargetURL, "proxy-pass", "", "HTTP反向代理URL")
	flag.IntVar(&opts.Port, "http-port", 8080, "端口号")
	flag.Parse()
}

func main() {
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", opts.Port))
	if err != nil {
		log.Fatalln(err)
	}

	handler, err := newHandler(opts.TargetURL)
	if err != nil {
		log.Fatalln(err)
	}

	if err = http.Serve(listener, handler); err != nil {
		log.Fatalln(err)
	}
}

type server struct {
	reverseHandler http.Handler // 反向代理处理器
	bizHandler     http.Handler // 业务处理器
}

func newHandler(rURL string) (http.Handler, error) {
	targetURL, err := url.Parse(rURL)
	if err != nil {
		return nil, err
	}

	reverseHandler := httputil.NewSingleHostReverseProxy(targetURL)

	bizHandler := http.NewServeMux()
	bizHandler.HandleFunc("/api/v1/ip", whatIsMyIP)

	return &server{
		reverseHandler: reverseHandler,
		bizHandler:     bizHandler,
	}, nil
}

func (server *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if uri := r.URL.RequestURI(); strings.HasPrefix(uri, "/api/v1/proxy/") {
		log.Printf("%s %s forwarded\n", r.Method, uri)
		server.reverseHandler.ServeHTTP(w, r)
	} else {
		log.Printf("Do request %s %s\n", r.Method, uri)
		server.bizHandler.ServeHTTP(w, r)
	}
}

func whatIsMyIP(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Fprintf(w, "Your IP is %s", ip)
}
