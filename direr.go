package tests

import (
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/beyondstorage/go-storage/v4/pairs"
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
				err := store.Delete(path, pairs.WithObjectMode(types.ModeDir))
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The first returning error should be nil", func() {
				So(err, ShouldBeNil)
			})

			o, err = d.CreateDir(path)
			Convey("The second returning error also should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object Mode should be dir", func() {
				// Dir object's mode must be Dir.
				So(o.Mode.IsDir(), ShouldBeTrue)
			})
		})

		Convey("When Create with ModeDir", func() {
			path := uuid.New().String()
			o := store.Create(path, pairs.WithObjectMode(types.ModeDir))

			defer func() {
				err := store.Delete(path, pairs.WithObjectMode(types.ModeDir))
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The Object Mode should be dir", func() {
				// Dir object's mode must be Dir.
				So(o.Mode.IsDir(), ShouldBeTrue)
			})
		})

		Convey("When Stat with ModeDir", func() {
			d, _ := store.(types.Direr)

			path := uuid.New().String()
			_, err := d.CreateDir(path)
			if err != nil {
				t.Error(err)
			}

			defer func() {
				err := store.Delete(path, pairs.WithObjectMode(types.ModeDir))
				if err != nil {
					t.Error(err)
				}
			}()

			o, err := store.Stat(path, pairs.WithObjectMode(types.ModeDir))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object Mode should be dir", func() {
				// Dir object's mode must be Dir.
				So(o.Mode.IsDir(), ShouldBeTrue)
			})
		})

		Convey("When Delete with ModeDir", func() {
			d, _ := store.(types.Direr)

			path := uuid.New().String()
			_, err := d.CreateDir(path)
			if err != nil {
				t.Error(err)
			}

			err = store.Delete(path, pairs.WithObjectMode(types.ModeDir))
			Convey("The first returning error should be nil", func() {
				So(err, ShouldBeNil)
			})

			err = store.Delete(path, pairs.WithObjectMode(types.ModeDir))
			Convey("The second returning error also should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
