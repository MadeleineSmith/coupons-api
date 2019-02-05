package handlers_test

import (
	"errors"
	"github.com/madeleinesmith/coupons/handlers"
	"github.com/madeleinesmith/coupons/handlers/handlersfakes"
	"github.com/madeleinesmith/coupons/model/coupon"
	"github.com/madeleinesmith/coupons/test_utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
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
			createdCoupon coupon.Coupon
			expectedResponse string
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
  "data": {
    "type": "coupons",
    "attributes": {
      "name": "Save £99 at Tesco",
      "brand": "Tesco",
      "value": 20
    }
  }
}`

			body := strings.NewReader(bodyJSON)
			request, err = http.NewRequest("POST", "/omg/lol", body)
			Expect(err).To(BeNil())

			name := "Save £99 at Tesco"
			brand := "Tesco"
			value := 20

			expectedCoupon = coupon.Coupon{
				Name: &name,
				Brand: &brand,
				Value: &value,
			}

			fakeCouponSerializer.DeserializeReturns(expectedCoupon, nil)

			createdCoupon = expectedCoupon
			createdCoupon.ID = "9dfd6d90-1c0a-11e9-9567-73937c5f9289"
			fakeCouponService.CreateCouponReturns(createdCoupon, nil)

			expectedResponse = `
{
  "data": {
    "type": "coupons",
    "id": "9dfd6d90-1c0a-11e9-9567-73937c5f9289",
    "attributes": {
      "name": "Save £99 at Tesco",
      "brand": "Tesco",
      "value": 20
    }
  }
}
`
			fakeCouponSerializer.SerializeReturns([]byte(expectedResponse), nil)
		})

		Context("Creating a coupon", func() {
			It("successfully creates a coupon", func() {
				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusCreated))
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Body.String()).To(Equal(expectedResponse))

				Expect(fakeCouponSerializer.DeserializeCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.DeserializeArgsForCall(0)).To(Equal([]byte(bodyJSON)))

				Expect(fakeCouponService.CreateCouponCallCount()).To(Equal(1))
				Expect(fakeCouponService.CreateCouponArgsForCall(0)).To(Equal(expectedCoupon))

				Expect(fakeCouponSerializer.SerializeCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.SerializeArgsForCall(0)).To(Equal(createdCoupon))
			})

			It("propagates the error if reading the request body fails", func() {
				request.Body = ioutil.NopCloser(test_utils.DummyReader{Message: "bad bad bad"})

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponSerializer.DeserializeCallCount()).To(Equal(0))
			})

			It("propagates the error if coupon deserialization fails", func() {
				fakeCouponSerializer.DeserializeReturns(coupon.Coupon{}, errors.New("good luck"))

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				Expect(fakeCouponSerializer.DeserializeCallCount()).To(Equal(1))
				Expect(fakeCouponService.CreateCouponCallCount()).To(Equal(0))
			})

			It("propagates the error if the coupon dbservice fails", func() {
				fakeCouponService.CreateCouponReturns(coupon.Coupon{}, errors.New("trololololol"))

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponService.CreateCouponCallCount()).To(Equal(1))
			})

			It("propagates the error if the coupon serialization fails", func() {
				fakeCouponSerializer.SerializeReturns(nil, errors.New("😱"))

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})

			// todo: move this test?
			It("errors if the method is unsupported", func() {
				request.Method = http.MethodDelete

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))

				Expect(fakeCouponSerializer.DeserializeCallCount()).To(Equal(0))
			})
		})
	})

	Describe("PATCH endpoint", func() {
		var (
			recorder *httptest.ResponseRecorder
			bodyJson string
			expectedCoupon coupon.Coupon
			request *http.Request
			handler handlers.CouponHandler
			fakeCouponSerializer *handlersfakes.FakeCouponSerializer
			fakeCouponService *handlersfakes.FakeCouponService
		)

		BeforeEach(func() {
			var err error

			fakeCouponService = &handlersfakes.FakeCouponService{}
			fakeCouponSerializer = &handlersfakes.FakeCouponSerializer{}

			handler = handlers.CouponHandler{
				Serializer: fakeCouponSerializer,
				CouponService: fakeCouponService,
			}

			bodyJson = `
					{
  "data": {
    "type": "coupons",
    "id": "0faec7ea-239f-11e9-9e44-d770694a0159",
    "attributes": {
      "brand": "Sainsbury's"
    }
  }
}`

			updateBody := strings.NewReader(bodyJson)
			recorder = httptest.NewRecorder()

			request, err = http.NewRequest("PATCH", "/omg/lol", updateBody)
			Expect(err).ToNot(HaveOccurred())

			brand := "Sainsbury's"

			expectedCoupon = coupon.Coupon{
				ID: "0faec7ea-239f-11e9-9e44-d770694a0159",
				Brand: &brand,
			}

			fakeCouponSerializer.DeserializeReturns(expectedCoupon, nil)
		})

		Context("Updating a coupon", func() {
			It("successfully updates a coupon", func() {
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))

				Expect(fakeCouponSerializer.DeserializeCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.DeserializeArgsForCall(0)).To(Equal([]byte(bodyJson)))

				Expect(fakeCouponService.UpdateCouponCallCount()).To(Equal(1))
				Expect(fakeCouponService.UpdateCouponArgsForCall(0)).To(Equal(expectedCoupon))
			})

			It("propagates the error if reading the request body fails", func() {
				request.Body = ioutil.NopCloser(test_utils.DummyReader{Message: "bad bad bad"})

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponSerializer.DeserializeCallCount()).To(Equal(0))
			})

			It("propagates the error if the coupon serializer fails", func() {
				fakeCouponSerializer.DeserializeReturns(coupon.Coupon{}, errors.New("Failed to deserialize to coupon instance"))

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponSerializer.DeserializeCallCount()).To(Equal(1))

				Expect(fakeCouponService.UpdateCouponCallCount()).To(Equal(0))
			})

			It("propagates the error if the db service fails", func() {
				fakeCouponService.UpdateCouponReturns(errors.New("db service failure"))

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponSerializer.DeserializeCallCount()).To(Equal(1))
				Expect(fakeCouponService.UpdateCouponCallCount()).To(Equal(1))
			})
		})
	})
})
