package main

import (
	"log"

	"github.com/CptAndy/hudsonsoftbackend/internal/db"
	"github.com/CptAndy/hudsonsoftbackend/internal/env"
	"github.com/CptAndy/hudsonsoftbackend/internal/store"
	"github.com/lpernett/godotenv"
	"go.uber.org/zap"
)

type Employee struct {
	ID       int64  `json:"id"`
	Emp_id   string `json:"emp_id"`
	Fname    string `json:"fname"`
	Lname    string `json:"lname"`
	Password string `json:"-"`
}

type Customer struct {
	ID    int64  `json:"id"`
	Fname string `json:"fname"`
	Lname string `json:"lname"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Product struct {
	ID             int64  `json:"id"`
	Product_name   string `json:"pname"`
	Sales_number   string `json:"snum"`
	Size           string `json:"size"`
	Price          string `json:"price"`
	Stock_quantity int64  `json:"stock"`
	Type_id        int64  `json:"Type"`
}

const version = ""

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/hudsonsoft?sslmode=disable"),
			maxOpenConns: env.GetInt("MAX_OPEN_CONNS", 20),
			maxIdleConns: env.GetInt("MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetString("MAX_IDLE_TIME", "15m"),
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime)
	if err != nil {
		logger.Panic(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))

}
