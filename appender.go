package tests

import (
	"bytes"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/beyondstorage/go-storage/v4/pairs"
	"github.com/beyondstorage/go-storage/v4/pkg/randbytes"
	"github.com/beyondstorage/go-storage/v4/types"
)

func TestAppender(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		ap, ok := store.(types.Appender)
		So(ok, ShouldBeTrue)

		Convey("When CreateAppend", func() {
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

		Convey("When CreateAppend with an existing appendable object", func() {
			path := uuid.NewString()
			o, err := ap.CreateAppend(path)

			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The first returning error should be nil", func() {
				So(err, ShouldBeNil)
			})

			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), size)

			_, err = ap.WriteAppend(o, r, size)
			if err != nil {
				t.Fatal(err)
			}

			o, err = ap.CreateAppend(path)

			Convey("The second returning error also should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object Mode should be appendable", func() {
				// Append object's mode must be appendable.
				So(o.Mode.IsAppend(), ShouldBeTrue)
			})

			Convey("The object append offset should be 0", func() {
				So(o.MustGetAppendOffset(), ShouldBeZeroValue)
			})
		})

		Convey("When Delete", func() {
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

			Convey("CommitAppend error should be nil", func() {
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
