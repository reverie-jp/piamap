#? help: ヘルプコマンド
help: Makefile
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^#?//p' $< | awk -F ':' '{ printf "  %-15s %s\n", $$1, $$2 }'
.PHONY: help

#? dev-up: 開発環境用サーバーを起動 (air で自動再ビルド)
dev-up:
	air
.PHONY: dev-up

#? proto: proto を Go / TS にコード生成
proto:
	cd proto && buf generate
.PHONY: proto

#? sqlc: SQL クエリを Go コードに変換
sqlc:
	sqlc generate
.PHONY: sqlc

#? grpcui: gRPC UI を起動
grpcui:
	grpcui -bind 0.0.0.0 -port 37611 -plaintext -open-browser=false localhost:50051
.PHONY: grpcui

#? migrate-up: データベースの構造をマイグレート
migrate-up:
	migrate -source file://migrations -database postgres://piamap:piamap@piamap-db:5432/piamap_db?sslmode=disable up
.PHONY: migrate-up

#? migrate-down: データベースの構造を初期化
migrate-down:
	migrate -source file://migrations -database postgres://piamap:piamap@piamap-db:5432/piamap_db?sslmode=disable down -all
.PHONY: migrate-down

#? psql: psql で DB に接続
psql:
	psql postgres://piamap:piamap@piamap-db:5432/piamap_db
.PHONY: psql
