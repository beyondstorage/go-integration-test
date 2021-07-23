package tests

import (
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/beyondstorage/go-storage/v4/pkg/randbytes"
	"github.com/beyondstorage/go-storage/v4/types"
)

func TestLinker(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		l, ok := store.(types.Linker)
		So(ok, ShouldBeTrue)

		Convey("When create a link object", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), size)
			target := uuid.New().String()

			_, err := store.Write(target, r, size)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(target)
				if err != nil {
					t.Error(err)
				}
			}()

			path := uuid.New().String()
			o, err := l.CreateLink(path, target)

			defer func() {
				err = store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The object mode should be link", func() {
				// Link object's mode must be link.
				So(o.Mode.IsLink(), ShouldBeTrue)
			})

			Convey("The linkTarget of the object must be the same as the target", func() {
				// The linkTarget must be the same as the target.
				linkTarget, ok := o.GetLinkTarget()

				So(ok, ShouldBeTrue)
				So(linkTarget, ShouldEqual, target)
			})

			Convey("Stat should get path object without error", func() {
				obj, err := store.Stat(path)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The object mode should be link", func() {
					// Link object's mode must be link.
					So(obj.Mode.IsLink(), ShouldBeTrue)
				})

				Convey("The linkTarget of the object must be the same as the target", func() {
					// The linkTarget must be the same as the target.
					linkTarget, ok := obj.GetLinkTarget()

					So(ok, ShouldBeTrue)
					So(linkTarget, ShouldEqual, target)
				})
			})
		})

		Convey("When create a link object from a not existing target", func() {
			target := uuid.New().String()

			path := uuid.New().String()
			o, err := l.CreateLink(path, target)

			defer func() {
				err = store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The object mode should be link", func() {
				// Link object's mode must be link.
				So(o.Mode.IsLink(), ShouldBeTrue)
			})

			Convey("The linkTarget of the object must be the same as the target", func() {
				linkTarget, ok := o.GetLinkTarget()

				So(ok, ShouldBeTrue)
				So(linkTarget, ShouldEqual, target)
			})

			Convey("Stat should get path object without error", func() {
				obj, err := store.Stat(path)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The object mode should be link", func() {
					// Link object's mode must be link.
					So(obj.Mode.IsLink(), ShouldBeTrue)
				})

				Convey("The linkTarget of the object must be the same as the target", func() {
					// The linkTarget must be the same as the target.
					linkTarget, ok := obj.GetLinkTarget()

					So(ok, ShouldBeTrue)
					So(linkTarget, ShouldEqual, target)
				})
			})

			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, _ := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))

			_, err = store.Write(target, bytes.NewReader(content), size)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(target)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("Read should get path object data without error", func() {
				var buf bytes.Buffer
				n, err := store.Read(path, &buf)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The content should be match", func() {
					So(n, ShouldNotBeNil)
					So(n, ShouldEqual, size)
					So(md5.Sum(buf.Bytes()), ShouldResemble, md5.Sum(content))
				})
			})
		})

		Convey("When CreateLink to an existing path", func() {
			firstSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), firstSize)
			firstTarget := uuid.New().String()

			_, err := store.Write(firstTarget, r, firstSize)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(firstTarget)
				if err != nil {
					t.Error(err)
				}
			}()

			path := uuid.New().String()
			o, err := l.CreateLink(path, firstTarget)

			defer func() {
				err = store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The first returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			secondSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			secondContent, _ := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), secondSize))
			secondTarget := uuid.New().String()

			_, err = store.Write(secondTarget, bytes.NewReader(secondContent), secondSize)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(secondTarget)
				if err != nil {
					t.Error(err)
				}
			}()

			o, err = l.CreateLink(path, secondTarget)

			Convey("The second returned error should also be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The object mode should be link", func() {
				// Link object's mode must be link.
				So(o.Mode.IsLink(), ShouldBeTrue)
			})

			Convey("The linkTarget of the object must be the same as the secondTarget", func() {
				// The linkTarget must be the same as the secondTarget.
				linkTarget, ok := o.GetLinkTarget()

				So(ok, ShouldBeTrue)
				So(linkTarget, ShouldEqual, secondTarget)
			})

			Convey("Read should get path object data without error", func() {
				var buf bytes.Buffer
				n, err := store.Read(path, &buf)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The content should be match", func() {
					// The content should match the secondTarget
					So(n, ShouldNotBeNil)
					So(n, ShouldEqual, secondSize)
					So(md5.Sum(buf.Bytes()), ShouldResemble, md5.Sum(secondContent))
				})
			})
		})
	})
}
