package handlers

import (
	"errors"
	"github.com/madeleinesmith/coupons/model/coupon"
	"io/ioutil"
	"net/http"
)

//go:generate counterfeiter . CouponService
type CouponService interface {
	CreateCoupon(couponInstance coupon.Coupon) (coupon.Coupon, error)
	UpdateCoupon(couponInstance coupon.Coupon) error
	GetCoupons() ([]*coupon.Coupon, error)
}

//go:generate counterfeiter . CouponSerializer
type CouponSerializer interface {
	DeserializeCoupon(bodyBytes []byte) (coupon.Coupon, error)
	SerializeCoupon(coupon coupon.Coupon) ([]byte, error)
	SerializeCoupons([]*coupon.Coupon) ([]byte, error)
}

type CouponHandler struct {
	Serializer CouponSerializer
	CouponService CouponService
}

func (h CouponHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		h.handlePost(w, req)
	} else if req.Method == http.MethodPatch {
		h.handlePatch(w, req)
	} else if req.Method == http.MethodGet {
		h.handleGet(w, req)
	} else {
		handleError(w, errors.New("Method not allowed"), http.StatusMethodNotAllowed)
	}
}

func (h CouponHandler) handlePost(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	couponInstance, err := h.Serializer.DeserializeCoupon(bodyBytes)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// TODO MS - add validation step here on couponInstance to assert that all fields are provided
	// otherwise return an error

	createdCoupon, err := h.CouponService.CreateCoupon(couponInstance)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	json, err := h.Serializer.SerializeCoupon(createdCoupon)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

func (h CouponHandler) handlePatch(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	couponInstance, err := h.Serializer.DeserializeCoupon(bodyBytes)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.CouponService.UpdateCoupon(couponInstance)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h CouponHandler) handleGet(w http.ResponseWriter, req *http.Request) {
	coupons, err := h.CouponService.GetCoupons()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	serializerCoupons, err := h.Serializer.SerializeCoupons(coupons)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serializerCoupons)
}

func handleError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}