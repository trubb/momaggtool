package main

import "context"

type Fund struct {
	AzaID int    `toml:"aza_id"`
	Name  string `toml:"name"`
}

type Funds struct {
	Funds []Fund `toml:"funds"`
}

type Config struct { // make this the core config instead?
	ctx           context.Context
	ctxCancel     context.CancelFunc
	fundsFilePath string
	Funds         []Fund
}
