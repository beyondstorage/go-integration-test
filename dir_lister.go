package tests

import (
	"io"
	"math/rand"
	"testing"

	ps "github.com/aos-dev/go-storage/v2/pairs"
	"github.com/aos-dev/go-storage/v2/pkg/randbytes"
	"github.com/aos-dev/go-storage/v2/types"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDirLister(t *testing.T, store types.Storager) {
	Convey("Given a dir lister", t, func() {
		var lister types.DirLister

		lister, ok := store.(types.DirLister)
		if !ok {
			t.Skip()
		}

		Convey("When List an empty dir", func() {
			it, err := lister.ListDir("")

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("The iterator should not be nil", func() {
				So(it, ShouldNotBeIn)
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
			_, err := store.Write(path, r, ps.WithSize(size))
			if err != nil {
				t.Error(err)
			}
			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			it, err := lister.ListDir("")
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("The iterator should not be nil", func() {
				So(it, ShouldNotBeIn)
			})

			o, err := it.Next()
			Convey("The name and size should be match", func() {
				So(o, ShouldNotBeNil)
				So(o.Name, ShouldEqual, path)

				osize, ok := o.GetSize()
				So(ok, ShouldBeTrue)
				So(osize, ShouldEqual, size)
			})
		})

	})
}
