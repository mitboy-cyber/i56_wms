module github.com/i56/i56-apps/i56-wms

go 1.25.0

require (
	github.com/i56/framework v1.1.0
	github.com/i56/modules v1.1.0
	golang.org/x/crypto v0.23.0
)

replace (
	github.com/i56/framework => ../../../i56-framework
	github.com/i56/modules => ../../modules
)
