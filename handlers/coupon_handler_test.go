package handlers_test

import (
	"github.com/coupons/handlers"
	"github.com/coupons/handlers/handlersfakes"
	"github.com/coupons/model/coupon"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Coupon Handler", func() {
	Describe("POST endpoint", func() {
		var (
			recorder *httptest.ResponseRecorder
			request *http.Request
			handler handlers.CouponHandler
			fakeCouponSerializer *handlersfakes.FakeCouponSerializer
			fakeCouponService *handlersfakes.FakeCouponService
			bodyJSON string
			expectedCoupon coupon.Coupon
		)

		BeforeEach(func() {
			var err error

			fakeCouponSerializer = &handlersfakes.FakeCouponSerializer{}
			fakeCouponService = &handlersfakes.FakeCouponService{}

			handler = handlers.CouponHandler{
				CouponService: fakeCouponService,
				Serializer: fakeCouponSerializer,
			}

			recorder = httptest.NewRecorder()

			bodyJSON = `{
	"name": "Save £99 at Tesco",
	"brand": "Tesco",
	"value": 20
}`

			updateBody := strings.NewReader(bodyJSON)
			request, err = http.NewRequest("POST", "/omg/lol", updateBody)
			Expect(err).To(BeNil())

			expectedCoupon = coupon.Coupon{
				Name: "Save £99 at Tesco",
				Brand: "Tesco",
				Value: 20,
			}

			fakeCouponSerializer.DeserializeReturns(expectedCoupon)
		})

		Context("Creating a coupon", func() {
			It("successfully creates a coupon", func() {
				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusCreated))

				Expect(fakeCouponSerializer.DeserializeCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.DeserializeArgsForCall(0)).To(Equal([]byte(bodyJSON)))

				Expect(fakeCouponService.CreateCouponCallCount()).To(Equal(1))
				Expect(fakeCouponService.CreateCouponArgsForCall(0)).To(Equal(expectedCoupon))
			})
		})
	})
})
