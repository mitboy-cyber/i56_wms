module github.com/i56/i56-apps/i56-wms

go 1.22

require (
	github.com/i56/framework v1.1.0
	github.com/i56/modules v1.1.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.10.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)

replace (
	github.com/i56/framework => ../../framework
	github.com/i56/modules => ../../modules
)
