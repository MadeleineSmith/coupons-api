package handlers

import (
	"github.com/coupons/model/coupon"
	"io/ioutil"
	"net/http"
)

//go:generate counterfeiter . CouponService
type CouponService interface {
	CreateCoupon(couponInstance coupon.Coupon)
}

//go:generate counterfeiter . CouponSerializer
type CouponSerializer interface {
	Deserialize(bodyBytes []byte) coupon.Coupon
}

type CouponHandler struct {
	CouponService CouponService
	Serializer CouponSerializer
}

func (h CouponHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handlePost(w, req)
}

func (h CouponHandler) handlePost(w http.ResponseWriter, req *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(req.Body)

	couponInstance := h.Serializer.Deserialize(bodyBytes)

	h.CouponService.CreateCoupon(couponInstance)

	w.WriteHeader(http.StatusCreated)
}