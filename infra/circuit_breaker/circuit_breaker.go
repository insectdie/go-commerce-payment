package circuitbreaker

import (
	"log"
	"time"

	"github.com/sony/gobreaker"
)

func NewCircuitBreakerInstance() *gobreaker.CircuitBreaker {
	st := gobreaker.Settings{
		Name: "integrationCircuitBreaker",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
		Timeout:     40 * time.Second,
		MaxRequests: 20,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			if to == gobreaker.StateOpen {
				log.Println("State Open!")
			}
			if from == gobreaker.StateOpen && to == gobreaker.StateHalfOpen {
				log.Println("Going from Open to Half Open")
			}
			if from == gobreaker.StateHalfOpen && to == gobreaker.StateClosed {
				log.Println("Going from Half Open to Closed!")
			}
		},
	}

	cb := gobreaker.NewCircuitBreaker(st)
	return cb
}

// MaxRequests:
// MaxRequests adalah jumlah maksimum permintaan yang diizinkan untuk melewati
// saat CircuitBreaker berada dalam status half-open.
// Jika MaxRequests bernilai 0, CircuitBreaker hanya mengizinkan 1 permintaan.

// Timeout:
// Timeout adalah periode ketika CircuitBreaker berada dalam status open,
// setelah itu status CircuitBreaker akan berubah menjadi half-open.
// Jika Timeout kurang dari atau sama dengan 0, nilai timeout CircuitBreaker akan diset menjadi 60 detik.

// ReadyToTrip:
// dijalankan setiap kali ada permintaan yang gagal saat CircuitBreaker berada di status closed.
// Jika ReadyToTrip mengembalikan nilai true, CircuitBreaker akan berubah ke status open.
// Jika ReadyToTrip tidak diatur (nil), akan digunakan aturan bawaan.
// Aturan bawaan akan mengubah CircuitBreaker ke status open jika terjadi lebih dari 5 kegagalan berturut-turut.
