package coupon_test

import (
	"github.com/madeleinesmith/coupons/model/coupon"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Coupon Serializer", func() {
	var s coupon.Serializer

	BeforeEach(func() {
		s = coupon.Serializer{}
	})

	Context("Deserializing", func() {
		It("deserializes a coupon", func() {

			bodyJSON := `{
  "data": {
    "type": "coupons",
    "attributes": {
      "name": "Save £99 at Tesco",
      "brand": "Tesco",
      "value": 20
    }
  }
}`
			body := []byte(bodyJSON)
			expectedCoupon := coupon.Coupon{
				Name: "Save £99 at Tesco",
				Brand: "Tesco",
				Value: 20,
			}

			s = coupon.Serializer{}
			model, err := s.Deserialize(body)
			Expect(err).To(Not(HaveOccurred()))

			Expect(model).To(Equal(expectedCoupon))
		})

		It("propagates the error", func() {
			s = coupon.Serializer{}

			_, err := s.Deserialize([]byte("🦄"))

			Expect(err).To(HaveOccurred())
		})
	})
})