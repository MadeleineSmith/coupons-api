package dbservices_test

import (
	"database/sql"
	"errors"
	"github.com/madeleinesmith/coupons/dbservices"
	"github.com/madeleinesmith/coupons/handlers"
	"github.com/madeleinesmith/coupons/model/coupon"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var _ = Describe("Coupon Service", func() {
	var (
		mockedService dbservices.CouponService
		dbMock        sqlmock.Sqlmock
		realService   dbservices.CouponService
	)

	BeforeEach(func() {
		var db *sql.DB
		var err error

		db, dbMock, err = sqlmock.New()
		Expect(err).NotTo(HaveOccurred())

		mockedService = dbservices.CouponService{
			DB: db,
		}

		realService = dbservices.CouponService{
			DB: realDB,
		}
	})

	Describe("CreateCoupon", func() {
		var (
			exampleCoupon coupon.Coupon
		)

		BeforeEach(func() {
			name := "Save £108 at Vue"
			brand := "Vue"
			value := 108

			exampleCoupon = coupon.Coupon{
				Name:  &name,
				Brand: &brand,
				Value: &value,
			}
		})

		It("successfully creates a coupon", func() {
			returnedCoupon, err := realService.CreateCoupon(exampleCoupon)
			Expect(err).ToNot(HaveOccurred())

			couponWithId := exampleCoupon
			couponWithId.ID = returnedCoupon.ID
			Expect(returnedCoupon).To(Equal(&couponWithId))

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

			_, err := mockedService.CreateCoupon(exampleCoupon)
			Expect(err).To(MatchError(ContainSubstring("oops I did it again 😇")))
			Expect(dbMock.ExpectationsWereMet()).To(Succeed())
		})
	})

	Describe("UpdateCoupon", func() {
		var (
			expectedCoupon coupon.Coupon
			updateQuery    string
		)

		BeforeEach(func() {
			brand := "Sainsbury's"
			name := "2 for 1 at Sainsbury's"
			value := 100

			expectedCoupon = coupon.Coupon{
				ID:    "0faec7ea-239f-11e9-9e44-d770694a0159",
				Name:  &name,
				Brand: &brand,
				Value: &value,
			}

			updateQuery = `UPDATE coupons SET name = \$1, brand = \$2, value = \$3 WHERE id = \$4`
		})

		It("successfully updates a coupon", func() {
			var newlyCreatedId string
			insertStatement := `INSERT INTO coupons (name, brand, value) VALUES ($1, $2, $3) RETURNING id`
			Expect(realDB.QueryRow(insertStatement, "A namely coupon", "Asda", 41).Scan(&newlyCreatedId)).To(Succeed())

			name := "A less namely coupon"
			value := 41

			couponToUpdate := coupon.Coupon{
				ID:    newlyCreatedId,
				Name:  &name,
				Value: &value,
			}

			Expect(realService.UpdateCoupon(couponToUpdate)).To(Succeed())

			capturedCoupon := coupon.Coupon{}
			Expect(realDB.QueryRow("SELECT name, brand, value FROM coupons WHERE id = $1", newlyCreatedId).Scan(&capturedCoupon.Name, &capturedCoupon.Brand, &capturedCoupon.Value)).To(Succeed())

			Expect(*capturedCoupon.Name).To(Equal(*couponToUpdate.Name))
			Expect(*capturedCoupon.Brand).To(Equal("Asda"))
			Expect(*capturedCoupon.Value).To(Equal(*couponToUpdate.Value))
		})

		It("propagates the error if exec fails", func() {
			dbMock.ExpectExec(updateQuery).
				WithArgs(*expectedCoupon.Name, *expectedCoupon.Brand, *expectedCoupon.Value, expectedCoupon.ID).
				WillReturnError(errors.New("oh dear 😭"))

			err := mockedService.UpdateCoupon(expectedCoupon)

			Expect(err).To(MatchError(ContainSubstring("oh dear 😭")))
			Expect(dbMock.ExpectationsWereMet()).To(Succeed())
		})
	})

	Describe("GetCoupons", func() {
		var (
			expectedCoupons []*coupon.Coupon
		)

		BeforeEach(func() {
			id1 := "354403f0-1c0e-11e9-9142-134e17ba9a5f"
			name1 := "Save £10 at Madeleine's Supermercado"
			brand1 := "Madeleine's"
			value1 := 10

			coupon1 := coupon.Coupon{
				ID:    id1,
				Name:  &name1,
				Brand: &brand1,
				Value: &value1,
			}

			id2 := "c614eeaa-1c9d-11e9-8c4f-3f7c43a05026"
			name2 := "Save £20 at Tom's Supermercado"
			brand2 := "Tom's"
			value2 := 20

			coupon2 := coupon.Coupon{
				ID:    id2,
				Name:  &name2,
				Brand: &brand2,
				Value: &value2,
			}

			id3 := "d56dab08-1c9d-11e9-944f-eb8190155a10"
			name3 := "Save £30 at Tom's Supermercado"
			brand3 := "Tom's"
			value3 := 30

			coupon3 := coupon.Coupon{
				ID:    id3,
				Name:  &name3,
				Brand: &brand3,
				Value: &value3,
			}

			expectedCoupons = []*coupon.Coupon{
				&coupon1,
				&coupon2,
				&coupon3,
			}

			_, err := realDB.Exec("INSERT INTO coupons (id, name, brand, value) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8), ($9, $10, $11, $12)",
				expectedCoupons[0].ID, *expectedCoupons[0].Name, *expectedCoupons[0].Brand, *expectedCoupons[0].Value,
				expectedCoupons[1].ID, *expectedCoupons[1].Name, *expectedCoupons[1].Brand, *expectedCoupons[1].Value,
				expectedCoupons[2].ID, *expectedCoupons[2].Name, *expectedCoupons[2].Brand, *expectedCoupons[2].Value)
			Expect(err).NotTo(HaveOccurred())
		})

		It("successfully retrieves coupons with no filter", func() {
			queryParams := handlers.Filters{}

			coupons, err := realService.GetCoupons(queryParams)
			Expect(err).NotTo(HaveOccurred())
			Expect(coupons).To(Equal(expectedCoupons))
		})

		It("successfully retrieves coupons with `brand` filter", func() {
			expectedBrand := "Tom's"

			coupons, err := realService.GetCoupons(handlers.Filters{
				Brand: &expectedBrand,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(coupons)).To(Equal(2))

			Expect(*coupons[0].Brand).To(Equal(expectedBrand))
			Expect(*coupons[1].Brand).To(Equal(expectedBrand))
		})

		It("successfully retrieves coupons with `value` filter", func() {
			expectedValue := 30

			coupons, err := realService.GetCoupons(handlers.Filters{
				Value: &expectedValue,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(coupons)).To(Equal(1))

			Expect(*coupons[0].Value).To(Equal(expectedValue))
		})

		It("successfully retrieves coupons with `name` filter", func() {
			expectedName := "Save £30 at Tom's Supermercado"

			coupons, err := realService.GetCoupons(handlers.Filters{
				Name: &expectedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(coupons)).To(Equal(1))

			Expect(*coupons[0].Name).To(Equal(expectedName))
		})

		It("successfully retrieves coupons with 2 (or more) filters", func() {
			expectedBrand := "Tom's"
			expectedValue := 30

			coupons, err := realService.GetCoupons(handlers.Filters{
				Brand: &expectedBrand,
				Value: &expectedValue,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(coupons)).To(Equal(1))

			Expect(*coupons[0].Brand).To(Equal("Tom's"))
			Expect(*coupons[0].Value).To(Equal(30))
		})

		It("propagates the error if querying the db fails", func() {
			dbMock.ExpectQuery("SELECT id, name, brand, value FROM coupons").WillReturnError(errors.New("boo 👻"))
			queryParams := handlers.Filters{}

			_, err := mockedService.GetCoupons(queryParams)
			Expect(err).To(MatchError("boo 👻"))

			Expect(dbMock.ExpectationsWereMet()).To(Succeed())
		})

		It("propagates the error if no rows are found", func() {
			queryParams := handlers.Filters{}

			rows := sqlmock.NewRows([]string{"id", "name", "brand", "value"})
			dbMock.ExpectQuery("SELECT id, name, brand, value FROM coupons").WillReturnRows(rows)

			_, err := mockedService.GetCoupons(queryParams)
			Expect(err).To(MatchError("sql: no rows in result set"))

			Expect(dbMock.ExpectationsWereMet()).To(Succeed())
		})

		It("propagates the error if scanning to the struct fails", func() {
			dbMock.ExpectQuery("SELECT id, name, brand, value FROM coupons").WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "brand", "value"}).
					AddRow(nil, nil, nil, nil))

			queryParams := handlers.Filters{}

			_, err := mockedService.GetCoupons(queryParams)
			Expect(err).To(HaveOccurred())

			Expect(dbMock.ExpectationsWereMet()).To(Succeed())
		})
	})

	Describe("GetCouponById", func() {
		It("successfully retrieves a coupon", func() {
			var couponId string

			insertStatement := `INSERT INTO coupons (name, brand, value) VALUES ($1, $2, $3) RETURNING id`
			Expect(realDB.QueryRow(insertStatement, "Save some money", "Accessorize", 10).Scan(&couponId)).To(Succeed())

			retrievedCoupon, err := realService.GetCouponById(couponId)
			Expect(err).ToNot(HaveOccurred())

			Expect(*retrievedCoupon.Name).To(Equal("Save some money"))
			Expect(*retrievedCoupon.Brand).To(Equal("Accessorize"))
			Expect(*retrievedCoupon.Value).To(Equal(10))
		})

		It("propagates the error if QueryRow/ scanning fails", func() {
			dbMock.ExpectQuery(`SELECT id, name, brand, value .*`).WillReturnError(sql.ErrNoRows)

			_, err := mockedService.GetCouponById("123")
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})
})
