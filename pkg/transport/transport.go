package transport

import (
	"fmt"

	"github.com/utkarsh-pro/use/pkg/storage"
	"github.com/utkarsh-pro/use/pkg/transport/http"
)

type TransportType string

type Transport interface {
	// Setup sets up the transport on the given address.
	Setup(addr string) error

	// Shutdown shuts down the transport.
	Shutdown() error
}

var (
	HTTPTransportType TransportType = "http"

	ErrInvalidTransportType = fmt.Errorf("invalid transport type")
)

// New returns a new transport
func New(transportType TransportType, storage storage.Storage) (Transport, error) {
	switch transportType {
	case HTTPTransportType:
		return http.New(storage), nil
	default:
		return nil, ErrInvalidTransportType
	}
}
