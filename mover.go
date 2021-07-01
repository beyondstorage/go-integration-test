package tests

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/beyondstorage/go-storage/v4/pkg/randbytes"
	"github.com/beyondstorage/go-storage/v4/services"
	"github.com/beyondstorage/go-storage/v4/types"
)

func TestMover(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		Convey("The Storager should implement Mover", func() {
			_, ok := store.(types.Mover)
			So(ok, ShouldBeTrue)
		})

		Convey("When Move a file", func() {
			m, _ := store.(types.Mover)

			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, _ := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			src := uuid.New().String()

			_, err := store.Write(src, bytes.NewReader(content), size)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(src)
				if err != nil {
					t.Error(err)
				}
			}()

			dst := uuid.New().String()
			err = m.Move(src, dst)

			defer func() {
				err = store.Delete(dst)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Stat should get src object not exist", func() {
				_, err := store.Stat(src)

				Convey("The error should be ErrObjectNotExist", func() {
					So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
				})
			})

			Convey("Read should get dst object data without error", func() {
				var buf bytes.Buffer
				n, err := store.Read(dst, &buf)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The content should be match", func() {
					So(buf, ShouldNotBeNil)
					So(n, ShouldEqual, size)
					So(md5.Sum(buf.Bytes()), ShouldResemble, md5.Sum(content))
				})
			})
		})
	})
}
