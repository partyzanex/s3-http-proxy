package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/partyzanex/s3-http-proxy/internal/endpoint"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Hostname  string
	EnableTLS bool

	MaxAvailableMemory  int
	MaxAvailableStorage int
	StoragePath         string

	S3Config endpoint.S3Config
}

func Run(ctx context.Context, config Config) error {
	e, err := endpoint.New(endpoint.Config{
		S3Config:    config.S3Config,
		MemorySize:  config.MaxAvailableMemory,
		StorageSize: config.MaxAvailableStorage,
		StoragePath: config.StoragePath,
	})
	if err != nil {
		return errors.Wrap(err, "cannot create endpoint")
	}

	s := fasthttp.Server{
		Handler: e.RequestHandler(),
	}

	go func() {
		select {
		case <-ctx.Done():
			err := s.Shutdown()
			if err != nil {
				log.Error(err)
			}
		}
	}()

	if err := s.ListenAndServe(config.Hostname); err != nil {
		return errors.Wrap(err, "ListenAndServe")
	}

	return nil
}
