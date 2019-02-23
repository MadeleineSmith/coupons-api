package validators_test

import (
	"github.com/madeleinesmith/coupons/model/coupon"
	"github.com/madeleinesmith/coupons/validators"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Coupon Validator", func() {
	var (
		couponValidator validators.CouponValidator
		sampleName      string
		sampleBrand     string
		sampleValue     int
		emptyField      string
	)

	BeforeEach(func() {
		couponValidator = validators.CouponValidator{}

		sampleName = "A Super Duper Coupon"
		sampleBrand = "Super Duper"
		sampleValue = 100
		emptyField = "     \n"
	})

	Context("With a valid coupon", func() {
		It("returns no error", func() {
			couponInstance := coupon.Coupon{
				Name:  &sampleName,
				Brand: &sampleBrand,
				Value: &sampleValue,
			}

			Expect(couponValidator.Validate(couponInstance)).To(Succeed())
		})
	})

	Context("With invalid fields", func() {
		DescribeTable("returns an error", func(coupon coupon.Coupon, errorMessage string) {
			err := couponValidator.Validate(coupon)
			Expect(err).To(HaveOccurred())

			Expect(err.Error()).To(Equal(errorMessage))
		},
			Entry("When the name field is not provided", coupon.Coupon{
				Name:  nil,
				Brand: &sampleBrand,
				Value: &sampleValue,
			}, "name field is required"),
			Entry("When the name field is empty", coupon.Coupon{
				Name:  &emptyField,
				Brand: &sampleBrand,
				Value: &sampleValue,
			}, "name field is required"),
			Entry("When the brand is not provided", coupon.Coupon{
				Name:  &sampleName,
				Brand: nil,
				Value: &sampleValue,
			}, "brand field is required"),
			Entry("When the brand field is empty", coupon.Coupon{
				Name:  &sampleName,
				Brand: &emptyField,
				Value: &sampleValue,
			}, "brand field is required"),
			Entry("When the value is not provided", coupon.Coupon{
				Name:  &sampleName,
				Brand: &sampleBrand,
				Value: nil,
			}, "value field is required"))
	})
})
