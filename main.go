package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var httpServer *http.Server
var loggerFile *os.File

func main() {

	AppInit()

	router := gin.New()
	rest := router.Group(config.Server.ContextPath)
	rest.POST("/schedule/absolute", AbsoluteSchedule)
	rest.POST("/schedule/relative", RelativeSchedule)
	rest.DELETE("/schedule/:id", AbortSchedule)

	httpServer = &http.Server{
		Addr:    ":" + IntToString(config.Server.Port),
		Handler: router,
	}

	go func() {
		// starts http server
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Application: Go-Trigger is UP")
	addShutDownHook(AppClose)
}

// AppInit - initialize application
func AppInit() {
	gin.SetMode(gin.ReleaseMode)
	// Add application up activities here.
	configInit()
	loggerInit()
}

func configInit() {
	configFile := flag.String("Config", "config.yaml", "configuration file path")
	LoadConfig(*configFile)
}

func loggerInit() {
	logPath := config.Server.LogPath
	loggerFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	multiLogger := io.MultiWriter(os.Stdout, loggerFile)
	log.SetOutput(multiLogger)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func addShutDownHook(shutDownHook func()) {
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignal
	shutDownHook()
}

// AppClose - shutdown hook
func AppClose() {
	log.Println("Application: Go-Trigger is shutting down")
	httpServerShutdown()
	// Add more shutdown activity here.
}

func httpServerShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
