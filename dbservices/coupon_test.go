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

var _ = Describe("Coupon Service", func() {
	Describe("CreateCoupon", func() {
		var (
			dbMock sqlmock.Sqlmock
			db *sql.DB
			couponService dbservices.CouponService
			insertQuery string
			exampleCoupon coupon.Coupon
		)

		BeforeEach(func() {
			var err error

			db, dbMock, err = sqlmock.New()
			Expect(err).ToNot(HaveOccurred())

			couponService = dbservices.CouponService{
				DB: db,
			}

			insertQuery = `INSERT INTO coupons \(name, brand, value\) VALUES \(\$1, \$2, \$3\)`
			exampleCoupon = coupon.Coupon{
				Name: "Save £108 at Vue",
				Brand: "Vue",
				Value: 108,
			}
		})

		It("successfully creates a coupon", func() {
			dbMock.ExpectExec(insertQuery).
				WithArgs(exampleCoupon.Name, exampleCoupon.Brand, exampleCoupon.Value).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := couponService.CreateCoupon(exampleCoupon)
			Expect(err).ToNot(HaveOccurred())
			Expect(dbMock.ExpectationsWereMet()).To(Succeed())
		})

		It("propagates the error", func() {
			dbMock.ExpectExec(insertQuery).
				WithArgs(exampleCoupon.Name, exampleCoupon.Brand, exampleCoupon.Value).
				WillReturnError(errors.New("oops I did it again 😇"))

			err := couponService.CreateCoupon(exampleCoupon)
			Expect(err).To(HaveOccurred())
			Expect(dbMock.ExpectationsWereMet()).To(Succeed())
		})
	})
})