package service

// Service is the interface for iCloud services
type Service interface {
	Start()
	Stop()
}
