package main

import (
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Post("POST /api/user/register", env.RegisterHandle)
}
