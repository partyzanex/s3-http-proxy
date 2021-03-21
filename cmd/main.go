package main

import (
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/partyzanex/s3-http-proxy/internal/endpoint"
	"github.com/partyzanex/s3-http-proxy/internal/server"

	"github.com/rs/zerolog/log"
)

func main() {
	app := cli.App{
		Usage:  "S3 caching proxy",
		Flags:  flags(),
		Before: before,
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err)
	}
}

func before(_ *cli.Context) error {
	r := http.NewServeMux()
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	go func() {
		err := http.ListenAndServe(":9081", r)
		if err != nil {
			log.Err(err)
		}
	}()

	return nil
}

func action(ctx *cli.Context) error {
	log.Printf("Listen %s", ctx.String("host"))

	err := server.Run(ctx.Context, server.Config{
		Hostname:            ctx.String("host"),
		EnableTLS:           ctx.Bool("tls"),
		MaxAvailableMemory:  ctx.Int("memory"),
		MaxAvailableStorage: ctx.Int("storage"),
		StoragePath:         ctx.String("storage-path"),
		S3Config: endpoint.S3Config{
			AccessKey:      ctx.String("s3-access-key"),
			SecretKey:      ctx.String("s3-secret-key"),
			Token:          ctx.String("s3-token"),
			Endpoint:       ctx.String("s3-endpoint"),
			Buckets:        ctx.StringSlice("s3-buckets"),
			Region:         ctx.String("s3-region"),
			DisableSSL:     ctx.Bool("s3-disable-ssl"),
			ForcePathStyle: ctx.Bool("s3-force-paths"),
		},
	})
	if err != nil {
		return errors.Wrap(err, "cannot run server")
	}

	return nil
}

func flags() []cli.Flag {
	const (
		Gigabyte = 1024 * 1024 * 1024
	)

	return []cli.Flag{
		&cli.StringFlag{
			Name:       "hostname",
			Aliases:    []string{"host", "address"},
			Usage:      "Server HTTP hostname (host)",
			EnvVars:    []string{"HOSTNAME"},
			Required:   true,
			Value:      "0.0.0.0:9080",
			HasBeenSet: true,
		},
		&cli.BoolFlag{
			Name:    "enable-tsl",
			Aliases: []string{"tls"},
			Usage:   "Switch TLS supports",
			EnvVars: []string{"ENABLE_TLS"},
		},
		&cli.IntFlag{
			Name:    "max-available-memory",
			Aliases: []string{"memory"},
			Usage:   "Max of available memory for cache",
			EnvVars: []string{"MAX_AVAILABLE_MEMORY"},
			Value:   Gigabyte,
		},
		&cli.IntFlag{
			Name:    "max-available-storage",
			Aliases: []string{"storage"},
			Usage:   "Max of available space in file storage",
			EnvVars: []string{"MAX_AVAILABLE_STORAGE"},
			Value:   Gigabyte,
		},
		&cli.StringFlag{
			Name:       "storage-path",
			EnvVars:    []string{"STORAGE_PATH"},
			Value:      "/tmp/s3-http-proxy-storage",
			HasBeenSet: true,
		},
		&cli.StringFlag{
			Name:    "s3-access-key",
			EnvVars: []string{"S3_ACCESS_KEY"},
		},
		&cli.StringFlag{
			Name:    "s3-secret-key",
			EnvVars: []string{"S3_SECRET_KEY"},
		},
		&cli.StringFlag{
			Name:    "s3-token",
			EnvVars: []string{"S3_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "s3-endpoint",
			EnvVars: []string{"S3_ENDPOINT"},
		},
		&cli.StringSliceFlag{
			Name:    "s3-buckets",
			EnvVars: []string{"S3_BUCKETS"},
		},
		&cli.StringFlag{
			Name:    "s3-region",
			EnvVars: []string{"S3_REGION"},
			Value:   "us-west-1",
		},
		&cli.BoolFlag{
			Name:    "s3-disable-ssl",
			EnvVars: []string{"S3_DISABLE_SSL"},
		},
		&cli.BoolFlag{
			Name:    "s3-force-paths",
			EnvVars: []string{"S3_FORCE_PATH_STYLE"},
		},
	}
}
