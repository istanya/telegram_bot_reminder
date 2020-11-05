-include app.env

migrateup:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose up

migratedown:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose down

sqlc:
	sqlc generate