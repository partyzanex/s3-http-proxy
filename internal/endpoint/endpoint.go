package endpoint

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/partyzanex/s3-http-proxy/internal/mime"
	"github.com/partyzanex/s3-http-proxy/pkg/storage"
	"github.com/partyzanex/s3-http-proxy/pkg/storage/file"
	"github.com/partyzanex/s3-http-proxy/pkg/storage/mem"
	"github.com/pkg/errors"
	"strings"
)

type Config struct {
	S3Config    S3Config
	MemorySize  int
	StorageSize int
	StoragePath string
}

type S3Config struct {
	AccessKey      string
	SecretKey      string
	Token          string
	Endpoint       string
	Buckets        []string
	Region         string
	DisableSSL     bool
	ForcePathStyle bool
}

type Endpoint struct {
	client       *s3.S3
	memory       storage.Interface
	files        storage.Interface
	mimeTypes    map[uint][]byte
	contentTypes map[string]uint
}

func New(config Config) (e *Endpoint, err error) {
	e = new(Endpoint)

	s3config := config.S3Config

	s3session, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			s3config.AccessKey, s3config.SecretKey, s3config.Token,
		),
		Endpoint:         aws.String(s3config.Endpoint),
		Region:           aws.String(s3config.Region),
		DisableSSL:       aws.Bool(s3config.DisableSSL),
		S3ForcePathStyle: aws.Bool(s3config.ForcePathStyle),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create a new session")
	}

	e.client = s3.New(s3session)

	e.files, err = file.New(config.StoragePath, config.StorageSize)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create a file storage")
	}

	mimeTypes, err := mime.Types()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get mime types")
	}

	e.memory = mem.New(config.MemorySize)
	e.mimeTypes = make(map[uint][]byte)
	e.contentTypes = make(map[string]uint)

	for i, mimeType := range mimeTypes {
		mimeType = strings.ToLower(mimeType)
		index := uint(i)
		e.mimeTypes[index] = []byte(mimeType)
		e.contentTypes[mimeType] = index
	}

	return
}
