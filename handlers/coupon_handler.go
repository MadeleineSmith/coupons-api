package handlers

import (
	"encoding/json"
	"github.com/coupons/model"
	"net/http"
)

type CouponService interface {
	CreateCoupon(coupon model.Coupon)
}

type CouponHandler struct {
	CouponService CouponService
}

func (h CouponHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handlePost(w, req)
}

func (h CouponHandler) handlePost(w http.ResponseWriter, req *http.Request) {
	// put following into function on a serializer
	var coupon model.Coupon

	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&coupon)

	// insert into db
	h.CouponService.CreateCoupon(coupon)

}