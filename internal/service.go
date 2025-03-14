package internal

import (
	"time"
)

const (
	TCP = "tcp"
)

type Service interface {
	HealthCheck(duration time.Duration) // Пишет в логи состояние сервиса
	Start() error
	Shutdown() error
	Configure() error
}
