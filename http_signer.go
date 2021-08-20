package tests

import (
	"bytes"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/beyondstorage/go-storage/v4/pkg/randbytes"
	"github.com/beyondstorage/go-storage/v4/types"
)

func TestHTTPSignerRead(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		signer, ok := store.(types.HTTPSigner)
		So(ok, ShouldBeTrue)

		Convey("When read via QuerySignHTTP", func() {
			size := rand.Int63n(4 * 1024 * 1024)
			content, err := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			path := uuid.New().String()
			_, err = store.Write(path, bytes.NewReader(content), size)
			if err != nil {
				t.Error(err)
			}
			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			req, err := signer.QuerySignHTTP(types.OpStoragerRead, path, time.Duration(time.Hour))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)

				So(req, ShouldNotBeNil)
				So(req.URL, ShouldNotBeNil)
			})

			client := http.Client{}
			resp, err := client.Do(req)
			Convey("The request returned error should be nil", func() {
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
			})

			defer resp.Body.Close()

			buf, err := ioutil.ReadAll(resp.Body)
			Convey("The content should be match", func() {
				So(err, ShouldBeNil)
				So(buf, ShouldNotBeNil)

				So(resp.ContentLength, ShouldEqual, size)
				So(sha256.Sum256(buf), ShouldResemble, sha256.Sum256(content))
			})
		})
	})
}

func TestHTTPSignerWrite(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		signer, ok := store.(types.HTTPSigner)
		So(ok, ShouldBeTrue)

		Convey("When write via QuerySignHTTP", func() {
			size := rand.Int63n(4 * 1024 * 1024)
			content, err := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			path := uuid.New().String()
			req, err := signer.QuerySignHTTP(types.OpStoragerWrite, path, time.Duration(time.Hour))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
				So(req, ShouldNotBeNil)
				So(req.URL, ShouldNotBeNil)
			})

			req.Body = ioutil.NopCloser(bytes.NewReader(content))
			req.ContentLength = size

			client := http.Client{}
			_, err = client.Do(req)
			Convey("The request returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("Read should get object data without error", func() {
				var buf bytes.Buffer
				n, err := store.Read(path, &buf)

				Convey("The content should be match", func() {
					So(err, ShouldBeNil)
					So(buf, ShouldNotBeNil)

					So(n, ShouldEqual, size)
					So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content))
				})
			})
		})
	})
}
