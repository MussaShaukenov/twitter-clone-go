package main

//func Serve(config *Config, routes http.Handler) error {
//	srv := &http.Server{
//		Addr:         config.addr,
//		Handler:      routes,
//		IdleTimeout:  time.Minute,
//		ReadTimeout:  10 * time.Second,
//		WriteTimeout: 10 * time.Second,
//	}
//
//	shutdownError := make(chan error)
//
//	go func() {
//		quit := make(chan os.Signal, 1)
//		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//		s := <-quit
//
//		config.logger.Infof("shutting down server, signal: %s", s.String())
//
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//
//		shutdownError <- srv.Shutdown(ctx)
//
//		os.Exit(0)
//	}()
//
//	config.logger.Info("starting server on port: ", srv.Addr)
//
//	err := srv.ListenAndServe()
//	if errors.Is(err, http.ErrServerClosed) {
//		return err
//	}
//
//	err = <-shutdownError
//	if err != nil {
//		return err
//	}
//
//	config.logger.Info("stopped server on addr: ", srv.Addr)
//
//	return nil
//}
