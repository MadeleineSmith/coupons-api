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
	"id": "0faec7ea-239f-11e9-9e44-d770694a0159",
    "attributes": {
      "name": "Save £99 at Tesco",
      "brand": "Tesco",
      "value": 20
    }
  }
}`
			body := []byte(bodyJSON)

			s = coupon.Serializer{}
			deserializeCoupon, err := s.Deserialize(body)
			Expect(err).To(Not(HaveOccurred()))

			Expect(deserializeCoupon.ID).To(Equal("0faec7ea-239f-11e9-9e44-d770694a0159"))
			Expect(*deserializeCoupon.Name).To(Equal("Save £99 at Tesco"))
			Expect(*deserializeCoupon.Brand).To(Equal("Tesco"))
			Expect(*deserializeCoupon.Value).To(Equal(20))
		})

		It("propagates the error", func() {
			s = coupon.Serializer{}

			_, err := s.Deserialize([]byte("🦄"))

			Expect(err).To(HaveOccurred())
		})
	})
})