package main

import (
	"context"
	"github.com/bamzi/jobrunner"
	"github.com/miguelsotocarlos/teleoma/internal/api/application"
	"github.com/miguelsotocarlos/teleoma/internal/api/config"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/check"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	database := db.NewServerDatabase(config.GetDatabase())

	err := database.Open()
	if err != nil {
		log.Print(err)
		return
	}

	defer check.Check(database.Close)
	crud.RefreshViews(database)
	// For jobs
	jobrunner.Start()

	engine := application.Setup(database)

	srv := &http.Server{
		Addr:    config.Port,
		Handler: engine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
