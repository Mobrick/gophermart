package ordersmigrations

import "embed"

//go:embed *.sql
var EmbedOrders embed.FS
