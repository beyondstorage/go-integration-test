package tests

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
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

			size := rand.Int63n(4 * 1024 * 1024)
			r := io.LimitReader(randbytes.NewRand(), size)
			rMD5Hex := calculateHashFromReader(r)
			src := uuid.New().String()

			_, err := store.Write(src, r, size)
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
					So(err, ShouldEqual, services.ErrObjectNotExist)
				})
			})

			Convey("Read should get dst object data without error", func() {
				var buf bytes.Buffer
				n, err := store.Read(dst, &buf)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The size should be equal", func() {
					So(n, ShouldEqual, size)
				})

				Convey("The hash should be equal", func() {
					bMD5 := md5.Sum(buf.Bytes())
					bMD5Hex := hex.EncodeToString(bMD5[:])
					So(rMD5Hex, ShouldEqual, bMD5Hex)
				})
			})
		})

		Convey("When Move to an existing file", func() {
			m, _ := store.(types.Mover)

			sSize := rand.Int63n(4 * 1024 * 1024)
			sr := io.LimitReader(randbytes.NewRand(), sSize)
			rMD5Hex := calculateHashFromReader(sr)
			src := uuid.New().String()

			_, err := store.Write(src, sr, sSize)
			if err != nil {
				t.Fatal(err)
			}

			dSize := rand.Int63n(4 * 1024 * 1024)
			dr := io.LimitReader(randbytes.NewRand(), dSize)
			dst := uuid.New().String()

			_, err = store.Write(dst, dr, dSize)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				sErr := store.Delete(src)
				dErr := store.Delete(dst)
				if sErr != nil {
					t.Error(sErr)
				}
				if dErr != nil {
					t.Error(dErr)
				}
			}()

			err = m.Move(src, dst)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Stat should get src object not exist", func() {
				_, err := store.Stat(src)

				Convey("The error should be ErrObjectNotExist", func() {
					So(err, ShouldEqual, services.ErrObjectNotExist)
				})
			})

			Convey("Read should get dst object data without error", func() {
				var buf bytes.Buffer
				n, err := store.Read(dst, &buf)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The size should be equal", func() {
					So(n, ShouldEqual, sSize)
				})

				Convey("The hash should be equal", func() {
					bMD5 := md5.Sum(buf.Bytes())
					bMD5Hex := hex.EncodeToString(bMD5[:])
					So(rMD5Hex, ShouldEqual, bMD5Hex)
				})
			})
		})

		Convey("When Move a non-existent file", func() {
			m, _ := store.(types.Mover)

			src := uuid.New().String()
			dst := uuid.New().String()

			err := m.Move(src, dst)
			Convey("The error should be ErrObjectNotExist", func() {
				So(err, ShouldEqual, services.ErrObjectNotExist)
			})
		})
	})
}
