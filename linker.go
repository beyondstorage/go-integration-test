package tests

import (
	"bytes"
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

func TestLinker(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		l, ok := store.(types.Linker)
		So(ok, ShouldBeTrue)

		Convey("When create a link object", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, _ := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			path := uuid.New().String()

			_, err := store.Write(path, bytes.NewReader(content), size)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			target := uuid.New().String()
			o, err := l.CreateLink(path, target)

			defer func() {
				err = store.Delete(target)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The object must should be link", func() {
				// Link object's mode must be link.
				So(o.Mode.IsLink(), ShouldBeTrue)
				// Link object's mode should not be read.
				So(o.Mode.IsRead(), ShouldBeFalse)
			})

			Convey("The linkTarget of the object must be the same as the target", func() {
				// The linkTarget must be the same as the target.
				linkTarget, ok := o.GetLinkTarget()

				So(ok, ShouldBeTrue)
				So(linkTarget, ShouldEqual, target)
			})
		})

		Convey("When create a link object from a not existing path", func() {
			path := uuid.New().String()

			target := uuid.New().String()

			Convey("Stat should get path object not exist", func() {
				// Path does not exist so there is no object on path
				_, err := store.Stat(path)

				Convey("The err should be ErrObjectNotExist", func() {
					So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
				})
			})

			o, err := l.CreateLink(path, target)

			Convey("The error should be ErrObjectNotExist", func() {
				So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
			})

			Convey("The object should be nil", func() {
				So(o, ShouldBeNil)
			})

			Convey("Stat should get target object mot exist", func() {
				// Path does not exist so no object will be created on target
				_, err = store.Stat(target)

				Convey("The err should be ErrObjectNotExist", func() {
					So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
				})
			})
		})

		Convey("When CreateLink to an existing target", func() {
			pathSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			pathContent, _ := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), pathSize))
			path := uuid.New().String()

			_, err := store.Write(path, bytes.NewReader(pathContent), pathSize)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			targetSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			targetContent, _ := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), targetSize))
			target := uuid.New().String()

			_, err = store.Write(target, bytes.NewReader(targetContent), targetSize)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(target)
				if err != nil {
					t.Error(err)
				}
			}()

			o, err := l.CreateLink(path, target)

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The object must should be link", func() {
				// Link object's mode must be link.
				So(o.Mode.IsLink(), ShouldBeTrue)
				// Link object's mode should not be read.
				So(o.Mode.IsRead(), ShouldBeFalse)
			})

			Convey("The linkTarget of the object must be the same as the target", func() {
				// The linkTarget must be the same as the target.
				linkTarget, ok := o.GetLinkTarget()

				So(ok, ShouldBeTrue)
				So(linkTarget, ShouldEqual, target)
			})

			Convey("List should get target ObjectIterator", func() {
				it, err := store.List(target)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The iterator should not be nil", func() {
					So(it, ShouldNotBeNil)
				})

				obj, err := it.Next()

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The object should not be nil", func() {
					So(obj, ShouldNotBeNil)
				})
			})
		})
	})
}
