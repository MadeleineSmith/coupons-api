package handlers_test

import (
	"errors"
	"github.com/lib/pq"
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
			recorder             *httptest.ResponseRecorder
			request              *http.Request
			handler              handlers.CouponHandler
			fakeCouponSerializer *handlersfakes.FakeCouponSerializer
			fakeCouponService    *handlersfakes.FakeCouponService
			bodyJSON             string
			expectedCoupon       coupon.Coupon
			createdCoupon        coupon.Coupon
			expectedResponse     string
		)

		BeforeEach(func() {
			var err error

			fakeCouponSerializer = &handlersfakes.FakeCouponSerializer{}
			fakeCouponService = &handlersfakes.FakeCouponService{}

			handler = handlers.CouponHandler{
				CouponService: fakeCouponService,
				Serializer:    fakeCouponSerializer,
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
				Name:  &name,
				Brand: &brand,
				Value: &value,
			}

			fakeCouponSerializer.DeserializeCouponReturns(expectedCoupon, nil)

			createdCoupon = expectedCoupon
			createdCoupon.ID = "9dfd6d90-1c0a-11e9-9567-73937c5f9289"
			fakeCouponService.CreateCouponReturns(&createdCoupon, nil)

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
			fakeCouponSerializer.SerializeCouponReturns([]byte(expectedResponse), nil)
		})

		Context("Creating a coupon", func() {
			It("successfully creates a coupon", func() {
				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusCreated))
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Body.String()).To(Equal(expectedResponse))

				Expect(fakeCouponSerializer.DeserializeCouponCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.DeserializeCouponArgsForCall(0)).To(Equal([]byte(bodyJSON)))

				Expect(fakeCouponService.CreateCouponCallCount()).To(Equal(1))
				Expect(fakeCouponService.CreateCouponArgsForCall(0)).To(Equal(expectedCoupon))

				Expect(fakeCouponSerializer.SerializeCouponCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.SerializeCouponArgsForCall(0)).To(Equal(&createdCoupon))
			})

			It("propagates the error if reading the request body fails", func() {
				request.Body = ioutil.NopCloser(test_utils.DummyReader{Message: "bad bad bad"})

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponSerializer.DeserializeCouponCallCount()).To(Equal(0))
			})

			It("propagates the error if coupon deserialization fails", func() {
				fakeCouponSerializer.DeserializeCouponReturns(coupon.Coupon{}, errors.New("good luck"))

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				Expect(fakeCouponSerializer.DeserializeCouponCallCount()).To(Equal(1))
				Expect(fakeCouponService.CreateCouponCallCount()).To(Equal(0))
			})

			It("propagates the error if the coupon dbservice fails", func() {
				fakeCouponService.CreateCouponReturns(nil, errors.New("trololololol"))

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponService.CreateCouponCallCount()).To(Equal(1))
			})

			It("propagates the error if the coupon serialization fails", func() {
				fakeCouponSerializer.SerializeCouponReturns(nil, errors.New("😱"))

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})

			// todo: move this test?
			It("errors if the method is unsupported", func() {
				request.Method = http.MethodDelete

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))

				Expect(fakeCouponSerializer.DeserializeCouponCallCount()).To(Equal(0))
			})
		})
	})

	Describe("PATCH endpoint", func() {
		var (
			recorder             *httptest.ResponseRecorder
			bodyJson             string
			expectedCoupon       coupon.Coupon
			request              *http.Request
			handler              handlers.CouponHandler
			fakeCouponSerializer *handlersfakes.FakeCouponSerializer
			fakeCouponService    *handlersfakes.FakeCouponService
		)

		BeforeEach(func() {
			var err error

			fakeCouponService = &handlersfakes.FakeCouponService{}
			fakeCouponSerializer = &handlersfakes.FakeCouponSerializer{}

			handler = handlers.CouponHandler{
				Serializer:    fakeCouponSerializer,
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
				ID:    "0faec7ea-239f-11e9-9e44-d770694a0159",
				Brand: &brand,
			}

			fakeCouponSerializer.DeserializeCouponReturns(expectedCoupon, nil)
		})

		Context("Updating a coupon", func() {
			It("successfully updates a coupon", func() {
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusNoContent))

				Expect(fakeCouponSerializer.DeserializeCouponCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.DeserializeCouponArgsForCall(0)).To(Equal([]byte(bodyJson)))

				Expect(fakeCouponService.UpdateCouponCallCount()).To(Equal(1))
				Expect(fakeCouponService.UpdateCouponArgsForCall(0)).To(Equal(expectedCoupon))
			})

			It("propagates the error if reading the request body fails", func() {
				request.Body = ioutil.NopCloser(test_utils.DummyReader{Message: "bad bad bad"})

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponSerializer.DeserializeCouponCallCount()).To(Equal(0))
			})

			It("propagates the error if the coupon serializer fails", func() {
				fakeCouponSerializer.DeserializeCouponReturns(coupon.Coupon{}, errors.New("Failed to deserialize to coupon instance"))

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponSerializer.DeserializeCouponCallCount()).To(Equal(1))

				Expect(fakeCouponService.UpdateCouponCallCount()).To(Equal(0))
			})

			It("propagates the error if the db service fails", func() {
				fakeCouponService.UpdateCouponReturns(errors.New("db service failure"))

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponSerializer.DeserializeCouponCallCount()).To(Equal(1))
				Expect(fakeCouponService.UpdateCouponCallCount()).To(Equal(1))
			})
		})
	})

	Describe("GET endpoint", func() {
		var (
			request              *http.Request
			fakeCouponSerializer handlersfakes.FakeCouponSerializer
			fakeCouponService    handlersfakes.FakeCouponService
			couponHandler        handlers.CouponHandler
			recorder             *httptest.ResponseRecorder
			couponsSlice         []*coupon.Coupon
		)

		BeforeEach(func() {
			var err error

			request, err = http.NewRequest(http.MethodGet, "/coupons", nil)
			Expect(err).NotTo(HaveOccurred())

			fakeCouponSerializer = handlersfakes.FakeCouponSerializer{}
			fakeCouponService = handlersfakes.FakeCouponService{}

			name1 := "2 for 1 at Vue"
			coupon1 := coupon.Coupon{
				ID:   "c1c16d12-1c0a-11e9-a3a3-9fd4e9cc6238",
				Name: &name1,
			}

			name2 := "Get 1000 Nectar points"
			coupon2 := coupon.Coupon{
				ID:   "f82df334-1c9b-11e9-afd2-070208c35e68",
				Name: &name2,
			}

			couponsSlice = []*coupon.Coupon{
				&coupon1,
				&coupon2,
			}

			fakeCouponService.GetCouponsReturns(couponsSlice, nil)

			fakeCouponSerializer.SerializeCouponsReturns([]byte(`😘`), nil)

			couponHandler = handlers.CouponHandler{
				Serializer:    &fakeCouponSerializer,
				CouponService: &fakeCouponService,
			}

			recorder = httptest.NewRecorder()
		})

		// /coupons
		Context("Getting coupons", func() {
			It("Successfully gets coupons", func() {
				couponHandler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(Equal(`😘`))
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

				Expect(fakeCouponService.GetCouponsCallCount()).To(Equal(1))

				Expect(fakeCouponSerializer.SerializeCouponsCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.SerializeCouponsArgsForCall(0)).To(Equal(couponsSlice))
			})

			It("propagates the error if the coupon service fails", func() {
				fakeCouponService.GetCouponsReturns(nil, errors.New("help 🙀"))

				couponHandler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponService.GetCouponsCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.SerializeCouponsCallCount()).To(Equal(0))
			})

			It("returns a 404 if no coupons are found", func() {
				fakeCouponService.GetCouponsReturns(nil, errors.New("sql: no rows in result set"))

				couponHandler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})

			It("returns a 400 if column is undefined", func() {
				errorPointer := new(pq.Error)
				errorPointer.Code = "42703"

				fakeCouponService.GetCouponsReturns(nil, errorPointer)

				couponHandler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("propagates the error if the coupon serializer fails", func() {
				fakeCouponSerializer.SerializeCouponsReturns(nil, errors.New("ahhhhh 😭"))

				couponHandler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponService.GetCouponsCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.SerializeCouponsCallCount()).To(Equal(1))
			})
		})

		// /coupons?brand=Madeleine's
		Context("Getting coupons with filter(s)", func() {
			It("Successfully retrieves coupons with 1 filter", func() {
				queryParameters := request.URL.Query()
				queryParameters.Add("brand", "Tom's")
				request.URL.RawQuery = queryParameters.Encode()

				couponHandler.ServeHTTP(recorder, request)

				Expect(fakeCouponService.GetCouponsCallCount()).To(Equal(1))
				// ConsistOf is required for variadic arguments as treats filters parameter as []handlers.Filter
				Expect(fakeCouponService.GetCouponsArgsForCall(0)).To(ConsistOf(handlers.Filter{
					FilterName:  "brand",
					FilterValue: "Tom's",
				}))
			})

			It("Successfully retrieves coupons with 2 (or more) filters", func() {
				queryParameters := request.URL.Query()
				queryParameters.Add("brand", "Tom's")
				queryParameters.Add("value", "30")
				request.URL.RawQuery = queryParameters.Encode()

				couponHandler.ServeHTTP(recorder, request)

				Expect(fakeCouponService.GetCouponsCallCount()).To(Equal(1))
				// ConsistOf is required for variadic arguments as treats filters parameter as []handlers.Filter
				Expect(fakeCouponService.GetCouponsArgsForCall(0)).To(ConsistOf(
					handlers.Filter{FilterName: "brand", FilterValue: "Tom's"},
					handlers.Filter{FilterName: "value", FilterValue: "30"}))
			})
		})
	})
})
