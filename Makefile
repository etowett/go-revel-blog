run:
	@revel run

build:
	@docker-compose build

up:
	@docker-compose up -d

logs:
	docker-compose logs -f

ps:
	@docker-compose ps

stop:
	@docker-compose stop

rm: stop
	@docker-compose rm

# make migration name=create_users
migration:
	@echo "Creating migration $(name)!"
	@goose -dir app/migrations create $(name) sql
	@echo "Done!"

migrate_up:
	@echo "Migrating up!"
	@goose -dir app/migrations postgres "host=127.0.0.1 port=5432 user=go-revel-blog password=go-revel-blog dbname=go-revel-blog sslmode=disable" up
	@echo "Done!"

migrate_down:
	@echo "Migrating down!"
	@goose -dir app/migrations postgres "host=127.0.0.1 port=5432 user=go-revel-blog password=go-revel-blog dbname=go-revel-blog sslmode=disable" down
	@echo "Done!"

migrate_status:
	@echo "Getting migration status!"
	@goose -dir app/migrations postgres "host=127.0.0.1 port=5432 user=go-revel-blog password=go-revel-blog dbname=go-revel-blog sslmode=disable" status
	@echo "Done!"

migrate_reset:
	@echo "Resetting migrations!"
	@goose -dir app/migrations postgres "host=127.0.0.1 port=5432 user=go-revel-blog password=go-revel-blog dbname=go-revel-blog sslmode=disable" reset
	@echo "Done!"

migrate_version:
	@echo "Getting migration version!"
	@goose -dir app/migrations postgres "host=127.0.0.1 port=5432 user=go-revel-blog password=go-revel-blog dbname=go-revel-blog sslmode=disable" version
	@echo "Done!"

migrate_redo: migrate_reset migrate_up
