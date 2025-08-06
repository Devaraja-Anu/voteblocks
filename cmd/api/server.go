package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

func (app *application) server() error {

	c := cron.New()

	_, err := c.AddFunc("@hourly", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := app.queries.DeactivateExpiredPolls(ctx)
		if err != nil {
			app.logger.PrintError(err, map[string]string{"cron": "deactivateExpiredPolls"})
		} else {
			app.logger.PrintInfo("Expired polls deactivated", nil)
		}
	})

	if err != nil {
		app.logger.PrintError(err, map[string]string{"cron": "failed to register cron job"})
	}

	serv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{"signal": s.String()})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := serv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.PrintInfo("completeing background tasks", map[string]string{
			"address": serv.Addr,
		})

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": serv.Addr,
	})

	err = serv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{"addr": serv.Addr})
	return nil
}
