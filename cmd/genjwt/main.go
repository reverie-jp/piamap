// genjwt: ローカル開発用に既存ユーザー (custom_id 指定) の access token を発行するユーティリティ。
// ユーザーが居なければ作成する。
//
//   docker compose exec -T piamap-api go run ./cmd/genjwt testuser01
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reverie-jp/piamap/internal/platform/jwt"
	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: genjwt <custom_id>")
		os.Exit(2)
	}
	customID := os.Args[1]
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://piamap:piamap@piamap-db:5432/piamap_db?sslmode=disable"
	}
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "db connect:", err)
		os.Exit(1)
	}
	defer db.Close()

	var idStr string
	err = db.QueryRow(ctx, "SELECT id FROM users WHERE custom_id = $1", customID).Scan(&idStr)
	var uid ulid.ULID
	if err == nil {
		uid, _ = ulid.Parse(idStr)
	} else if errors.Is(err, pgx.ErrNoRows) {
		uid = ulid.New()
		if _, err := db.Exec(ctx,
			"INSERT INTO users (id, custom_id, display_name) VALUES ($1, $2, $3)",
			uid, customID, customID); err != nil {
			fmt.Fprintln(os.Stderr, "user insert:", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stderr, "user lookup:", err)
		os.Exit(1)
	}

	secret := os.Getenv("AUTH_JWT_SECRET_KEY")
	if secret == "" {
		fmt.Fprintln(os.Stderr, "AUTH_JWT_SECRET_KEY is not set")
		os.Exit(1)
	}
	jm := jwt.NewManager(secret, 24*time.Hour, 30*24*time.Hour)
	access, err := jm.GenerateAccessToken(uid)
	if err != nil {
		fmt.Fprintln(os.Stderr, "generate token:", err)
		os.Exit(1)
	}
	fmt.Println(access)
}
