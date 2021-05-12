package tests

import (
	"bytes"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/aos-dev/go-storage/v3/pairs"
	"github.com/aos-dev/go-storage/v3/pkg/randbytes"
	"github.com/aos-dev/go-storage/v3/types"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAppender(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		Convey("The Storager should implement Appender", func() {
			_, ok := store.(types.Appender)
			So(ok, ShouldBeTrue)
		})

		Convey("When CreateAppend", func() {
			ap, _ := store.(types.Appender)

			path := uuid.NewString()
			o, err := ap.CreateAppend(path)

			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object Mode should be appendable", func() {
				// Append object's mode must be appendable.
				So(o.Mode.IsAppend(), ShouldBeTrue)
			})
		})

		Convey("When Delete", func() {
			ap, _ := store.(types.Appender)

			path := uuid.NewString()
			_, err := ap.CreateAppend(path)
			if err != nil {
				t.Error(err)
			}

			err = store.Delete(path)
			Convey("The first returning error should be nil", func() {
				So(err, ShouldBeNil)
			})

			err = store.Delete(path)
			Convey("The second returning error also should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When WriteAppend", func() {
			ap, _ := store.(types.Appender)

			path := uuid.NewString()
			o, err := ap.CreateAppend(path)
			if err != nil {
				t.Error(err)
			}

			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, _ := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			r := bytes.NewReader(content)

			n, err := ap.WriteAppend(o, r, size)

			Convey("WriteAppend error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("WriteAppend size should be equal to n", func() {
				So(n, ShouldEqual, size)
			})
		})

		Convey("When CommitAppend", func() {
			ap, _ := store.(types.Appender)

			path := uuid.NewString()
			o, err := ap.CreateAppend(path)
			if err != nil {
				t.Error(err)
			}

			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, _ := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			r := bytes.NewReader(content)

			_, err = ap.WriteAppend(o, r, size)
			if err != nil {
				t.Error(err)
			}

			err = ap.CommitAppend(o)

			Convey("CommitAppend error should be nil", func(){
				So(err, ShouldBeNil)
			})

			var buf bytes.Buffer
			_, err = store.Read(path, &buf, pairs.WithSize(size))

			Convey("Read error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("The content should be match", func() {
				So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content))
			})
		})
	})

}
