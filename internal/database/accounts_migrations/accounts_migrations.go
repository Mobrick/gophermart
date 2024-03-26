package accountsmigrations

import "embed"

//go:embed *.sql
var EmbedAccounts embed.FS
