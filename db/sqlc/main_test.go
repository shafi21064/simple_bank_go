package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shafi21064/simplebank/util"
)

var testQuaries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannor load config: ", err)
	}
	testDB, err = pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("can't connct db: ", err)
	}
	testQuaries = New(testDB)
	os.Exit(m.Run())
}
