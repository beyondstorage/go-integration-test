package tests

import (
	"io"
	"math/rand"
	"testing"

	"github.com/aos-dev/go-storage/v3/pairs"
	"github.com/aos-dev/go-storage/v3/pkg/randbytes"
	"github.com/aos-dev/go-storage/v3/types"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMultiparter(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		Convey("The Storager should implement Multiparter", func() {
			_, ok := store.(types.Multiparter)
			So(ok, ShouldBeTrue)
		})

		Convey("When CreateMultipart", func() {
			m, _ := store.(types.Multiparter)

			path := uuid.New().String()
			o, err := m.CreateMultipart(path)

			defer func() {
				err := store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object Mode should be part", func() {
				// Multipart object's mode must be Part.
				So(o.Mode.IsPart(), ShouldBeTrue)
				// Multipart object's mode must not be Read.
				So(o.Mode.IsRead(), ShouldBeFalse)
			})

			Convey("The Object must have multipart id", func() {
				// Multipart object must have multipart id.
				_, ok := o.GetMultipartID()
				So(ok, ShouldBeTrue)
			})
		})

		Convey("When Delete with multipart id", func() {
			m, _ := store.(types.Multiparter)

			path := uuid.New().String()
			o, err := m.CreateMultipart(path)
			if err != nil {
				t.Error(err)
			}

			err = store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When Delete with multipart id twice", func() {
			m, _ := store.(types.Multiparter)

			path := uuid.New().String()
			o, err := m.CreateMultipart(path)
			if err != nil {
				t.Error(err)
			}

			for i := 0; i < 2; i++ {
				err = store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})
			}
		})

		Convey("When Stat with multipart id", func() {
			m, _ := store.(types.Multiparter)

			path := uuid.New().String()
			o, err := m.CreateMultipart(path)
			if err != nil {
				t.Error(err)
			}

			multipartId := o.MustGetMultipartID()

			defer func() {
				err := store.Delete(path, pairs.WithMultipartID(multipartId))
				if err != nil {
					t.Error(err)
				}
			}()

			mo, err := store.Stat(path, pairs.WithMultipartID(multipartId))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
				So(mo, ShouldNotBeNil)
			})

			Convey("The Object Mode should be part", func() {
				// Multipart object's mode must be Part.
				So(mo.Mode.IsPart(), ShouldBeTrue)
				// Multipart object's mode must not be Read.
				So(mo.Mode.IsRead(), ShouldBeFalse)
			})

			Convey("The Object must have multipart id", func() {
				// Multipart object must have multipart id.
				mid, ok := mo.GetMultipartID()
				So(ok, ShouldBeTrue)
				So(mid, ShouldEqual, multipartId)
			})
		})

		Convey("When Create with multipart id", func() {
			m, _ := store.(types.Multiparter)

			path := uuid.New().String()
			o, err := m.CreateMultipart(path)
			if err != nil {
				t.Error(err)
			}

			multipartId := o.MustGetMultipartID()

			defer func() {
				err := store.Delete(path, pairs.WithMultipartID(multipartId))
				if err != nil {
					t.Error(err)
				}
			}()

			mo := store.Create(path, pairs.WithMultipartID(multipartId))

			Convey("The Object Mode should be part", func() {
				// Multipart object's mode must be Part.
				So(mo.Mode.IsPart(), ShouldBeTrue)
				// Multipart object's mode must not be Read.
				So(mo.Mode.IsRead(), ShouldBeFalse)
			})

			Convey("The Object must have multipart id", func() {
				// Multipart object must have multipart id.
				mid, ok := mo.GetMultipartID()
				So(ok, ShouldBeTrue)
				So(mid, ShouldEqual, multipartId)
			})
		})

		Convey("When WriteMultipart", func() {
			m, _ := store.(types.Multiparter)

			path := uuid.New().String()
			o, err := m.CreateMultipart(path)
			if err != nil {
				t.Error(err)
			}

			defer func() {
				err := store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
				if err != nil {
					t.Error(err)
				}
			}()

			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), size)

			n, err := m.WriteMultipart(o, r, size, 0)

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The size should be match", func() {
				So(n, ShouldEqual, size)
			})
		})

		Convey("When ListMultiPart", func() {
			m, _ := store.(types.Multiparter)

			path := uuid.New().String()
			o, err := m.CreateMultipart(path)
			if err != nil {
				t.Error(err)
			}

			defer func() {
				err := store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
				if err != nil {
					t.Error(err)
				}
			}()

			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			partNumber := rand.Intn(1000)        // Choose a random part number from [0, 1000)
			r := io.LimitReader(randbytes.NewRand(), size)

			_, err = m.WriteMultipart(o, r, size, partNumber)
			if err != nil {
				t.Error(err)
			}

			it, err := m.ListMultipart(o)

			Convey("ListMultipart error should be nil", func() {
				So(err, ShouldBeNil)
				So(it, ShouldNotBeNil)
			})

			p, err := it.Next()
			Convey("Next error should be nil", func() {
				So(err, ShouldBeNil)
				So(it, ShouldNotBeNil)
			})
			Convey("The part number and size should be match", func() {
				So(p.Index, ShouldEqual, partNumber)
				So(p.Size, ShouldEqual, size)
			})
		})

		Convey("When CompletePart", func() {
			m, _ := store.(types.Multiparter)

			path := uuid.New().String()
			o, err := m.CreateMultipart(path)
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
			partNumber := rand.Intn(1000)        // Choose a random part number from [0, 1000)
			r := io.LimitReader(randbytes.NewRand(), size)

			_, err = m.WriteMultipart(o, r, size, partNumber)
			if err != nil {
				t.Error(err)
			}

			err = m.CompleteMultipart(o, []*types.Part{
				{Index: partNumber, Size: size},
			})

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The object should be readable after complete", func() {
				ro, err := store.Stat(path)

				So(err, ShouldBeNil)
				So(ro.Mode.IsRead(), ShouldBeTrue)
				So(ro.Mode.IsPart(), ShouldBeFalse)
			})
		})
	})
}
