package scheduler

import (
	"libraryOnline/services"
	"log"
	"time"
)

func StartReservationScheduler(service *services.ReservationService) {
	go func() {
		for {
			if err := service.CancelExpired(); err != nil {
				log.Println("Error cancelando reservas expiradas:", err)
			}
			time.Sleep(1 * time.Hour) // Revisa cada hora
		}
	}()
}
