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

	Context("Serializing", func() {
		It("serializes a coupon", func() {
			name := "Save £20 at Madz supermarkets"
			brand := "Madz supermarkets"
			value := 20

			exampleCoupon := coupon.Coupon{
				ID: "658a191a-28b5-11e9-9968-87c211c8c951",
				Name: &name,
				Brand: &brand,
				Value: &value,
			}

			serializer := coupon.Serializer{}

			byteSlice, _ := serializer.Serialize(exampleCoupon)

			// Find better way to do the below
			Expect(string(byteSlice)).To(Equal(`{"data":{"type":"coupons","id":"658a191a-28b5-11e9-9968-87c211c8c951","attributes":{"brand":"Madz supermarkets","name":"Save £20 at Madz supermarkets","value":20}}}
`))
		})
	})
})