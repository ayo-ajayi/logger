package main

import (
	"fmt"
	"time"

	"github.com/ayo-ajayi/logger"
)

func myLogCallback(event logger.EventInfo) {
	fmt.Printf("Log at %v, file %s, line %d: %s\n", event.Time(), event.File(), event.Line(), event.Message())

}
func main() {
	log := logger.NewLogger(logger.INFO, false)
	log.AddCallback(myLogCallback)
	log.SetUseColor(false)

	port := 6565
	if port < 1 || port > 65535 {
		log.Error("Port must be a number between 1 and 65535.")
		return
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	page := fmt.Sprintf("{\"name\": \"John Doe\", \"age\": 30, \"date\": \"%s\"}", currentTime)

	log.Info("Response page: %s", page)
	log.Error("Server started on port %d", port)
}
