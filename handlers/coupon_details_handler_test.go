package handlers_test

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/madeleinesmith/coupons/handlers"
	"github.com/madeleinesmith/coupons/handlers/handlersfakes"
	"github.com/madeleinesmith/coupons/model/coupon"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("CouponDetailsHandler", func() {
	Describe("GET endpoint", func() {
		Context("Getting a coupon", func() {
			var (
				request *http.Request
				recorder *httptest.ResponseRecorder
				couponId string
				fakeCouponService handlersfakes.FakeCouponService
				fakeCouponSerializer handlersfakes.FakeCouponSerializer
				handler handlers.CouponDetailsHandler
				sampleCoupon *coupon.Coupon
			)

			BeforeEach(func() {
				var err error

				request, err = http.NewRequest(http.MethodGet, "/omg/lol", nil)
				Expect(err).ToNot(HaveOccurred())

				recorder = httptest.NewRecorder()

				couponId = "123"

				queryParams := map[string]string{
					"couponId": couponId,
				}

				request = mux.SetURLVars(request, queryParams)

				fakeCouponService = handlersfakes.FakeCouponService{}
				fakeCouponSerializer = handlersfakes.FakeCouponSerializer{}

				sampleCoupon = new(coupon.Coupon)
				fakeCouponService.GetCouponByFilterReturns(sampleCoupon, nil)
				fakeCouponSerializer.SerializeCouponReturns([]byte("halfway there 🙏"), nil)

				handler = handlers.CouponDetailsHandler{
					CouponService: &fakeCouponService,
					Serializer: &fakeCouponSerializer,
				}
			})

			It("successfully retrieves a coupon", func() {
				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

				Expect(fakeCouponService.GetCouponByFilterCallCount()).To(Equal(1))
				filterName, filterValue := fakeCouponService.GetCouponByFilterArgsForCall(0)
				Expect(filterName).To(Equal("id"))
				Expect(filterValue).To(Equal(couponId))

				Expect(fakeCouponSerializer.SerializeCouponCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.SerializeCouponArgsForCall(0)).To(Equal(sampleCoupon))

				Expect(string(recorder.Body.Bytes())).To(Equal("halfway there 🙏"))
			})

			It("errors if the couponId URL variable is not set", func() {
				var emptyURLVars map[string]string
				request = mux.SetURLVars(request, emptyURLVars)

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

				Expect(fakeCouponService.GetCouponByFilterCallCount()).To(Equal(0))
				Expect(fakeCouponSerializer.SerializeCouponCallCount()).To(Equal(0))

				Expect(string(recorder.Body.Bytes())).To(ContainSubstring("couponId URL variable not found"))
			})

			It("returns a 404 if the coupon id does not exist", func() {
				fakeCouponService.GetCouponByFilterReturns(&coupon.Coupon{}, sql.ErrNoRows)

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusNotFound))

				Expect(fakeCouponService.GetCouponByFilterCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.SerializeCouponCallCount()).To(Equal(0))

				Expect(string(recorder.Body.Bytes())).To(ContainSubstring("sql: no rows in result set"))
			})

			It("propagates the error if the db service fails", func() {
				fakeCouponService.GetCouponByFilterReturns(&coupon.Coupon{}, errors.New("🎷🎷🎷🎷"))

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponService.GetCouponByFilterCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.SerializeCouponCallCount()).To(Equal(0))

				Expect(string(recorder.Body.Bytes())).To(ContainSubstring("🎷🎷🎷🎷"))
			})

			It("propagates the error if the coupon serializer fails", func() {
				fakeCouponSerializer.SerializeCouponReturns([]byte(""), errors.New("shocking 👻"))

				handler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))

				Expect(fakeCouponService.GetCouponByFilterCallCount()).To(Equal(1))
				Expect(fakeCouponSerializer.SerializeCouponCallCount()).To(Equal(1))

				Expect(string(recorder.Body.Bytes())).To(ContainSubstring("shocking 👻"))
			})

			// probs doesn't belong in the GET context but cba to do all the setup all over again
			It("errors if the http method is not supported", func() {
				request.Method = http.MethodOptions

				handler.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
				Expect(string(recorder.Body.Bytes())).To(ContainSubstring("Method not allowed"))
			})
		})
	})
})