package main

import (
	_ "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity"
	"github.com/yon-module/yon-framework/server"
)

func main() {
	s := server.NewServer()
	s.Start()
}
