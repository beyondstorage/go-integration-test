package tests

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"

	ps "github.com/aos-dev/go-storage/v3/pairs"
	"github.com/aos-dev/go-storage/v3/pkg/randbytes"
	"github.com/aos-dev/go-storage/v3/services"
	"github.com/aos-dev/go-storage/v3/types"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStorager(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		var err error

		Convey("The Storager should not be nil", func() {
			So(store, ShouldNotBeNil)
		})

		Convey("The error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("When String called", func() {
			s := store.String()

			Convey("The string should not be empty", func() {
				So(s, ShouldNotBeEmpty)
			})
		})

		Convey("When Metadata called", func() {
			m, err := store.Metadata()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The metadata should not be empty", func() {
				So(m, ShouldNotBeEmpty)
			})
		})

		Convey("When Read a file", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
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

			var buf bytes.Buffer

			n, err := store.Read(path, &buf)

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The content should be match", func() {
				So(buf, ShouldNotBeNil)

				So(n, ShouldEqual, size)
				So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content))
			})

		})

		Convey("When Write a file", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), size)
			path := uuid.New().String()

			_, err := store.Write(path, r, size)

			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Stat should get Object without error", func() {
				o, err := store.Stat(path)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The name and size should be match", func() {
					So(o, ShouldNotBeNil)
					So(o.Path, ShouldEqual, path)

					osize, ok := o.GetContentLength()
					So(ok, ShouldBeTrue)
					So(osize, ShouldEqual, size)
				})
			})

			Convey("Read should get Object data without error", func() {
				var buf bytes.Buffer
				n, err := store.Read(path, &buf)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The size should be equal", func() {
					So(n, ShouldEqual, size)
				})
			})

		})

		Convey("When Stat a file", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
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

			o, err := store.Stat(path)

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object name and size should be match", func() {
				So(o, ShouldNotBeNil)
				So(o.Path, ShouldEqual, path)

				osize, ok := o.GetContentLength()
				So(ok, ShouldBeTrue)
				So(osize, ShouldEqual, size)
			})
		})

		Convey("When Delete a file", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, err := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			path := uuid.New().String()
			_, err = store.Write(path, bytes.NewReader(content), size)
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

			Convey("Stat should get nil Object and ObjectNotFound error", func() {
				o, err := store.Stat(path)

				So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
				So(o, ShouldBeNil)
			})
		})

		Convey("When List an empty dir", func() {
			it, err := store.List("", ps.WithListMode(types.ListModeDir))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("The iterator should not be nil", func() {
				So(it, ShouldNotBeNil)
			})

			o, err := it.Next()

			Convey("The next should be done", func() {
				So(err, ShouldBeError, types.IterateDone)
			})
			Convey("The object should be nil", func() {
				So(o, ShouldBeNil)
			})
		})

		Convey("When List a dir within files", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), size)
			path := uuid.New().String()
			_, err := store.Write(path, r, size)
			if err != nil {
				t.Error(err)
			}
			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			it, err := store.List("",  ps.WithListMode(types.ListModeDir))
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("The iterator should not be nil", func() {
				So(it, ShouldNotBeNil)
			})

			o, err := it.Next()
			Convey("The name and size should be match", func() {
				So(o, ShouldNotBeNil)
				So(o.Path, ShouldEqual, path)

				osize, ok := o.GetContentLength()
				So(ok, ShouldBeTrue)
				So(osize, ShouldEqual, size)
			})
		})
	})
}
