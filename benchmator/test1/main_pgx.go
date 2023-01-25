package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kochanac/coursework2023/benchmator/pkg/benchmator"
	"github.com/jackc/pgx/v5/pgxpool"
)

func establishConn(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(pgDSN)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 1
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	return db, err
}

const ROWNUM int64 = 2e7
const TIMEOUT = time.Second * 10
const RPS = 15
const DURATION = time.Minute * 5

const sumrequests = RPS * DURATION / time.Second

func main() {
	db, err := establishConn(context.Background())
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	oneSelect := func() {
		ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
		defer cancel()

		item := getRandFromSeq()

		row := db.QueryRow(ctx,
			fmt.Sprintf("select id, data from %s where id=$1", pgTable), item.id)

		var itemGot test1
		err = row.Scan(&itemGot.id, &itemGot.data)
		if err != nil {
			log.Printf("%s", err)
		}

		if itemGot.data == 0 {
			log.Printf("%v %v", itemGot, item)
			return
		}
		if itemGot.data != item.data {
			log.Fatalf("not matching id=%d: %s != %s", item.id, item.data, itemGot.data)
		}
	}

	stats := benchmator.Run(oneSelect, DURATION, RPS)

	log.Println(stats)

	wr, err := os.OpenFile("stats.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	err = stats.Save(wr)
	if err != nil {
		log.Fatalf("failed to save to a file: %s", err)
	}
}
