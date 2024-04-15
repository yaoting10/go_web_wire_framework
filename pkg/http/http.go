package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, addr string) {

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func SendGetRequest(url string, headers map[string]string) string {
	var req *http.Request
	req, _ = http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	return string(respByte)
}

func SendProxyGet(proxy string, target string) string {
	var client http.Client
	if proxy != "" {
		proxyUrl, _ := url.Parse(proxy)
		client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, Timeout: 30 * time.Second}
	} else {
		client = http.Client{Timeout: 30 * time.Second}
	}

	resp, err := client.Get(target)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	return string(respByte)
}
