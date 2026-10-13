package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aaksant/school-api/internal/api/classes"
	"github.com/aaksant/school-api/internal/api/middlewares"
	"github.com/aaksant/school-api/internal/api/refresh_tokens"
	"github.com/aaksant/school-api/internal/api/students"
	"github.com/aaksant/school-api/internal/api/teachers"
	"github.com/aaksant/school-api/internal/api/teaching_assignments"
	"github.com/aaksant/school-api/internal/api/users"
	"github.com/aaksant/school-api/internal/auth"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

const port string = "3000"
const (
	issuer          = "school-api"
	accessTokenTtl  = 5 * time.Minute
	refreshTokenTtl = 7 * 24 * time.Hour
)

const maxBodyBytes = 1 << 20 // 1 MiB

func main() {
	conn, err := sql.Open("sqlite", "db/index.db")
	if err != nil {
		log.Fatalln("error on sql.Open:", err)
	}
	defer conn.Close()

	mux := http.NewServeMux()

	if err := godotenv.Load(); err != nil {
		log.Fatalln("no .env file found")
	}

	key := os.Getenv("JWT_SECRET")
	if key == "" {
		log.Fatalln("invalid JWT secret")
	}

	tokenManager, err := auth.NewTokenManager(
		[]byte(key),
		accessTokenTtl,
		issuer,
	)
	if err != nil {
		log.Fatalln("error creating token manager")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	studentsSqliteRepo := students.NewSqliteRepository(conn)
	studentsHandler := students.NewHandler(studentsSqliteRepo, validate)
	studentsHandler.RegisterRoutes(mux, tokenManager)

	teachersSqliteRepo := teachers.NewSqliteRepository(conn)
	teachersHandler := teachers.NewHandler(teachersSqliteRepo, validate)
	teachersHandler.RegisterRoutes(mux, tokenManager)

	classesSqliteRepo := classes.NewSqliteRepository(conn)
	classesHandler := classes.NewHandler(classesSqliteRepo, validate)
	classesHandler.RegisterRoutes(mux, tokenManager)

	teachingAssignmentSqliteRepo := teaching_assignments.NewSqliteRepository(conn)
	teachingAssignmentsHandler := teaching_assignments.NewHandler(
		teachingAssignmentSqliteRepo,
		validate,
	)
	teachingAssignmentsHandler.RegisterRoutes(mux, tokenManager)

	usersSqliteRepo := users.NewSqliteRepository(conn)
	refreshTokensRepo := refresh_tokens.NewSqliteRepository(conn)

	usersHandler := users.NewHandler(
		usersSqliteRepo,
		refreshTokensRepo,
		validate,
		tokenManager,
		refreshTokenTtl,
	)
	usersHandler.RegisterRoutes(mux, tokenManager)

	rl := middlewares.NewIpRateLimiter(10, 20)
	hpp := middlewares.Hpp(middlewares.HppOptions{
		CheckQueryParams:        true,
		CheckBodyParams:         true,
		CheckBodyForContentType: "application/x-www-form-urlencoded",
		Whitelist:               []string{"tag"}, // temporary, example
	})
	maxBytes := middlewares.MaxBytes(maxBodyBytes)

	handler := middlewares.Chain(middlewares.ChainConfig{
		Handler: mux,
		Middlewares: []middlewares.Middleware{
			maxBytes,
			middlewares.Cors,
			rl.RateLimiter,
			middlewares.SecurityHeaders,
			middlewares.Compression,
			hpp,
		},
	})

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		TLSConfig:         tlsConfig,
	}

	// graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Println("server is running on port:", port)
		if err := server.ListenAndServeTLS(
			"certs/cert.pem",
			"certs/key.pem",
		); err != nil && err != http.ErrServerClosed {
			log.Fatalln("error starting the server:", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalln("forced shutdown:", err)
	}

	log.Println("server exited gracefully")
}
