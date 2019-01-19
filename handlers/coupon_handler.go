package handlers

import (
	"github.com/coupons/model/coupon"
	"io"
	"net/http"
)

type CouponService interface {
	CreateCoupon(coupon coupon.Coupon)
}

type CouponSerializer interface {
	Deserialize(body io.ReadCloser) coupon.Coupon
}

type CouponHandler struct {
	CouponService CouponService
	Serializer CouponSerializer
}

func (h CouponHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handlePost(w, req)
}

func (h CouponHandler) handlePost(w http.ResponseWriter, req *http.Request) {
	coupon := h.Serializer.Deserialize(req.Body)

	h.CouponService.CreateCoupon(coupon)
}