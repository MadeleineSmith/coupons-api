package handlers

import (
	"github.com/coupons/model/coupon"
	"io/ioutil"
	"net/http"
)

//go:generate counterfeiter . CouponService
type CouponService interface {
	CreateCoupon(couponInstance coupon.Coupon) error
}

//go:generate counterfeiter . CouponSerializer
type CouponSerializer interface {
	Deserialize(bodyBytes []byte) (coupon.Coupon, error)
}

type CouponHandler struct {
	CouponService CouponService
	Serializer CouponSerializer
}

func (h CouponHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handlePost(w, req)
}

func (h CouponHandler) handlePost(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	couponInstance, err := h.Serializer.Deserialize(bodyBytes)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	err = h.CouponService.CreateCoupon(couponInstance)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// consider placing elsewhere
func handleError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}