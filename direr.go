package tests

import (
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/beyondstorage/go-storage/v4/types"
)

func TestDirer(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		Convey("The Storager should implement Direr", func() {
			_, ok := store.(types.Direr)
			So(ok, ShouldBeTrue)
		})

		Convey("When CreateDir", func() {
			d, _ := store.(types.Direr)

			path := uuid.New().String()
			o, err := d.CreateDir(path)

			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object Mode should be dir", func() {
				// Dir object's mode must be Dir.
				So(o.Mode.IsDir(), ShouldBeTrue)
			})
		})
	})
}
