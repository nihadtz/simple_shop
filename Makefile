.PHONY: migrate downgrade

migrate:
	migrate -source "file://migrations" -database "mysql://root:test@tcp(127.0.0.1:4306)/shop" up 1

downgrade:
	migrate -source "file://migrations" -database "mysql://root:test@tcp(127.0.0.1:4306)/shop" down 1