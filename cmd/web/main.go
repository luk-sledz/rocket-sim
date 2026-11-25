package main

import (
	"log"
	"rocket-sim/internal/server"
)

func main() {
	// ListenAndServe blokuje: proces żyje, dopóki serwer działa (lub padnie błąd).
	if err := server.ListenAndServe(":8080", "web"); err != nil {
		log.Fatal(err) // wypisuje błąd i zakończa z kodem ≠ 0
	}
}
