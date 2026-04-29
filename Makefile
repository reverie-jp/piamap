#? help: ヘルプコマンド
help: Makefile
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^#?//p' $< | awk -F ':' '{ printf "  %-15s %s\n", $$1, $$2 }'
.PHONY: help

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
