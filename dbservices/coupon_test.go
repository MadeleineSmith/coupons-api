package dbservices_test

import (
	"database/sql"
	"errors"
	"github.com/madeleinesmith/coupons/dbservices"
	"github.com/madeleinesmith/coupons/model/coupon"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// TODO: figure out how to test with a real database
// Test happy path with real db

var _ = Describe("Coupon Service", func() {
	Describe("CreateCoupon", func() {
		var (
			dbMock sqlmock.Sqlmock
			db *sql.DB
			couponService dbservices.CouponService
			exampleCoupon coupon.Coupon
		)

		BeforeEach(func() {
			var err error

			db, dbMock, err = sqlmock.New()
			Expect(err).ToNot(HaveOccurred())

			couponService = dbservices.CouponService{
				DB: db,
			}

			name := "Save £108 at Vue"
			brand := "Vue"
			value := 108

			exampleCoupon = coupon.Coupon{
				Name: &name,
				Brand: &brand,
				Value: &value,
			}
		})

		// use real db for happy path
		// use mocked db for error cases
		// always use real db for integration tests

		It("successfully creates a coupon", func() {
			couponService = dbservices.CouponService{
				DB: realDB,
			}

			returnedCoupon, err := couponService.CreateCoupon(exampleCoupon)
			Expect(err).ToNot(HaveOccurred())

			couponWithId := exampleCoupon
			couponWithId.ID = returnedCoupon.ID
			Expect(returnedCoupon).To(Equal(couponWithId))

			var capturedCoupon coupon.Coupon

			Expect(realDB.QueryRow("SELECT id, name, brand, value FROM coupons WHERE id=$1", returnedCoupon.ID).
				Scan(&capturedCoupon.ID, &capturedCoupon.Name, &capturedCoupon.Brand, &capturedCoupon.Value)).To(Succeed())

			Expect(capturedCoupon.ID).NotTo(BeEmpty())
			Expect(*capturedCoupon.Name).To(Equal("Save £108 at Vue"))
			Expect(*capturedCoupon.Brand).To(Equal("Vue"))
			Expect(*capturedCoupon.Value).To(Equal(108))
		})

		It("propagates the error", func() {
			dbMock.ExpectQuery("INSERT INTO coupons .*").
				WillReturnError(errors.New("oops I did it again 😇"))

			_, err := couponService.CreateCoupon(exampleCoupon)
			Expect(err).To(MatchError(ContainSubstring("oops I did it again 😇")))
			Expect(dbMock.ExpectationsWereMet()).To(Succeed())
		})
	})

	Describe("UpdateCoupon", func() {
		var (
			dbMock sqlmock.Sqlmock
			couponService dbservices.CouponService
			expectedCoupon coupon.Coupon
			updateQuery string
		)

		BeforeEach(func() {
			var db *sql.DB
			var err error

			db, dbMock, err = sqlmock.New()
			Expect(err).ToNot(HaveOccurred())

			couponService = dbservices.CouponService{
				DB: db,
			}

			brand := "Sainsbury's"
			name := "2 for 1 at Sainsbury's"
			value := 100

			expectedCoupon = coupon.Coupon{
				ID: "0faec7ea-239f-11e9-9e44-d770694a0159",
				Name: &name,
				Brand: &brand,
				Value: &value,
			}

			updateQuery = `UPDATE coupons SET name = \$1, brand = \$2, value = \$3 WHERE id = \$4`
		})

		It("successfully updates a coupon", func() {
			dbMock.ExpectExec(updateQuery).
				WithArgs(*expectedCoupon.Name, *expectedCoupon.Brand, *expectedCoupon.Value, expectedCoupon.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := couponService.UpdateCoupon(expectedCoupon)
			Expect(err).ToNot(HaveOccurred())
			Expect(dbMock.ExpectationsWereMet()).To(Succeed())
		})

		It("propagates the error if exec fails", func() {
			dbMock.ExpectExec(updateQuery).
				WithArgs(*expectedCoupon.Name, *expectedCoupon.Brand, *expectedCoupon.Value, expectedCoupon.ID).
				WillReturnError(errors.New("oh dear 😭"))

			err := couponService.UpdateCoupon(expectedCoupon)

			Expect(err).To(MatchError(ContainSubstring("oh dear 😭")))
			Expect(dbMock.ExpectationsWereMet()).To(Succeed())
		})
	})
})