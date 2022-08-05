package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"withdrawal-exchanger/internal/app/worker"
	"withdrawal-exchanger/internal/clients/infura"
	"withdrawal-exchanger/internal/connectors/logger"
	"withdrawal-exchanger/internal/connectors/postgres"
	"withdrawal-exchanger/internal/pkg/deposit_events/repository"
	"withdrawal-exchanger/internal/pkg/deposit_events/usecase"
	"withdrawal-exchanger/internal/utils/env"
)

// stages:
//
// 1. prepare patrons for threads
// 2. FetchData
// 3. FilterData by WC
// 4. Exchange pubKey to validator index
// 5. Store into DB
// https://dune.com/queries/567094/1062483

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, envErr := env.Read(ctx)
	if envErr != nil {
		println("Read env error:", envErr.Error())

		os.Exit(1)
	}

	conn, err := infura.NewHttp(ctx, cfg.AppConfig.InfuraKey)
	if err != nil {
		println("Could not establish connection to infura!", err.Error())

		os.Exit(1)
	}

	db, errDB := postgres.Connect(cfg.PgConfig)
	if errDB != nil {
		println("Connect db error:", errDB.Error())
		os.Exit(1)
	}

	f, err := os.OpenFile("./app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	logger := logger.New(f)

	repo := repository.New(db)

	uc := usecase.New(conn,
		repo,
		cfg.AppConfig.DepositContractAddress,
		cfg.AppConfig.WC,
		cfg.AppConfig.InfuraProjectID,
		cfg.AppConfig.InfuraProjectSecret,
		cfg.AppConfig.InfuraETH2Host,
		&http.Client{},
	)

	w := worker.New(uc, logger)

	start := time.Now()

	w.Run(ctx, w.PrepareChPatrons(cfg.AppConfig.FromBlock, cfg.AppConfig.ToBlock, cfg.AppConfig.Step)...)

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}
