//go:build tools
// +build tools

//go:generate go build -o ../bin/abigen github.com/ethereum/go-ethereum/cmd/abigen
//go:generate go build -tags 'postgres' -o ../bin/migrate github.com/golang-migrate/migrate/v4/cmd/migrate


// Package tools contains go:generate commands for all project tools with versions stored in local go.mod file
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/ethereum/go-ethereum/cmd/abigen"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)