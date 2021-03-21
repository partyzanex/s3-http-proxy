package endpoint

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/partyzanex/s3-http-proxy/pkg/pipeline"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/partyzanex/s3-http-proxy/pkg/storage"
)

type FileGetter func(ctx *fasthttp.RequestCtx, bucket, key []byte) error

func (e *Endpoint) getFileGetter() FileGetter {
	var pathSeparator = []byte("/")

	return func(ctx *fasthttp.RequestCtx, bucket, key []byte) error {
		var (
			reader      io.ReadCloser
			contentType []byte
			storageKey  = storage.NewKeyFromBytes(bytes.Join([][]byte{bucket, key}, pathSeparator))
		)

		defer func() {
			if reader != nil {
				_ = reader.Close()
			}
		}()

		err := pipeline.Run(
			func() (next bool, err error) {
				obj, err := e.memory.Get(storageKey)
				if err != nil {
					return true, errors.Wrap(err, "cannot get object from memory")
				}

				reader = io.NopCloser(obj.Reader())
				contentType = obj.GetMimeType(e.mimeTypes)

				return false, nil
			},
			func() (next bool, err error) {
				obj, err := e.files.Get(storageKey)
				if err != nil {
					return true, errors.Wrap(err, "cannot get object from files")
				}

				err = e.memory.Set(storageKey, obj)
				if err != nil {
					return false, errors.Wrap(err, "cannot set object to memory")
				}

				reader = io.NopCloser(obj.Reader())
				contentType = obj.GetMimeType(e.mimeTypes)

				return false, nil
			},
			func() (next bool, err error) {
				output, err := e.client.GetObject(&s3.GetObjectInput{
					Bucket: aws.String(string(bucket)),
					Key:    aws.String(string(key)),
				})
				if err != nil {
					return true, errors.Wrap(err, "cannot get object from s3")
				}

				obj := e.newObjectFromS3(output)

				err = e.files.Set(storageKey, obj)
				if err != nil {
					return false, errors.Wrap(err, "cannot set object to files")
				}

				reader = io.NopCloser(obj.Reader())
				contentType = obj.GetMimeType(e.mimeTypes)

				return false, nil
			},
		)
		if err != nil {
			return errors.Wrap(err, "cannot get object")
		}

		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.Header.SetBytesV(fasthttp.HeaderContentType, contentType)

		_, err = io.Copy(ctx, reader)
		if err != nil {
			return errors.Wrap(err, "cannot copy object")
		}

		return nil
	}
}

func (e *Endpoint) newObjectFromS3(output *s3.GetObjectOutput) *storage.Object {
	b, _ := ioutil.ReadAll(output.Body)
	contentType := strings.ToLower(*output.ContentType)
	_ = output.Body.Close()

	return &storage.Object{
		Body:         b,
		ContentType:  e.contentTypes[contentType],
		LastModified: time.Now().UnixNano(),
	}
}
