package app

import (
	"avito_internship/internal/repository"
	"avito_internship/internal/transport"
)

func Run() {
	repository.Connect()
	transport.Run()
}
