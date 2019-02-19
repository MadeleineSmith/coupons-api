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

	Context("DeserializeCoupon", func() {
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
			deserializeCoupon, err := s.DeserializeCoupon(body)
			Expect(err).To(Not(HaveOccurred()))

			Expect(deserializeCoupon.ID).To(Equal("0faec7ea-239f-11e9-9e44-d770694a0159"))
			Expect(*deserializeCoupon.Name).To(Equal("Save £99 at Tesco"))
			Expect(*deserializeCoupon.Brand).To(Equal("Tesco"))
			Expect(*deserializeCoupon.Value).To(Equal(20))
		})

		It("propagates the error", func() {
			s = coupon.Serializer{}

			_, err := s.DeserializeCoupon([]byte("🦄"))

			Expect(err).To(HaveOccurred())
		})
	})

	Context("SerializeCoupon", func() {
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

			byteSlice, _ := serializer.SerializeCoupon(&exampleCoupon)

			Expect(string(byteSlice)).To(MatchJSON(`{
   "data":{
      "type":"coupons",
      "id":"658a191a-28b5-11e9-9968-87c211c8c951",
      "attributes":{
         "brand":"Madz supermarkets",
         "name":"Save £20 at Madz supermarkets",
         "value":20
      }
   }
}`))
		})
	})

	Context("SerializeCoupons", func() {
		It("serializes multiple coupons", func() {
			id1 := "354403f0-1c0e-11e9-9142-134e17ba9a5f"
			name1 := "Save £10 at Madeleine's Supermercado"
			brand1 := "Madeleine's"
			value1 := 10

			coupon1 := coupon.Coupon{
				ID: id1,
				Name: &name1,
				Brand: &brand1,
				Value: &value1,
			}

			id2 := "c614eeaa-1c9d-11e9-8c4f-3f7c43a05026"
			name2 := "Save £20 at Tom's Supermercado"
			brand2 := "Tom's"
			value2 := 20

			coupon2 := coupon.Coupon{
				ID: id2,
				Name: &name2,
				Brand: &brand2,
				Value: &value2,
			}

			expectedCoupons := []*coupon.Coupon{
				&coupon1,
				&coupon2,
			}

			byteSlice, err := s.SerializeCoupons(expectedCoupons)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(byteSlice)).To(MatchJSON(`{
  "data":[
    {
      "type": "coupons",
      "id": "354403f0-1c0e-11e9-9142-134e17ba9a5f",
      "attributes": {
        "brand": "Madeleine's",
        "name": "Save £10 at Madeleine's Supermercado",
        "value": 10
      }
    },
    {
      "type": "coupons",
      "id": "c614eeaa-1c9d-11e9-8c4f-3f7c43a05026",
      "attributes": {
        "brand": "Tom's",
        "name": "Save £20 at Tom's Supermercado",
        "value": 20
      }
    }
  ]
}`))
		})
	})
})