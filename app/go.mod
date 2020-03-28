module github.com/agundry/rps_financial/app

go 1.13

require (
	github.com/agundry/rps_financial/app/models v0.0.0-20200328204548-87fb24729743
	github.com/agundry/rps_financial/util v0.1.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/mux v1.7.4
)

replace github.com/agundry/rps_financial/app/models => ./models

replace github.com/agundry/rps_financial/util => ../util
