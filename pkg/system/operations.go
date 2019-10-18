package system

import (
	"context"
	"errors"
	"time"
)

// ErrNotImplemented declares error for method that isn't implemented
var ErrNotImplemented = errors.New("this method is not implemented")

// ErrEmptyServerPointer declares error for nil pointer
var ErrEmptyServerPointer = errors.New("server pointer should not be nil")

// Operations implements simplest Operator interface
type Operations struct {
	gracePeriod time.Duration
	shutdowns   []Shutdowner
}

// NewOperator creates operator
func NewOperator(cfg *Config, sd ...Shutdowner) *Operations {
	service := new(Operations)
	service.gracePeriod = time.Duration(cfg.Grace.Period) * time.Second
	service.shutdowns = append(service.shutdowns, sd...)

	return service
}

// Reload operation implementation
func (o Operations) Reload() error {
	return ErrNotImplemented
}

// Maintenance operation implementation
func (o Operations) Maintenance() error {
	return ErrNotImplemented
}

// Shutdown operation
func (o Operations) Shutdown() []error {
	var errs []error

	ctx, cancel := context.WithTimeout(context.TODO(), o.gracePeriod)

	defer cancel()

	for _, fn := range o.shutdowns {
		if err := fn.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
