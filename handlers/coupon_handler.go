package handlers

import (
	"errors"
	"github.com/madeleinesmith/coupons/model/coupon"
	"io/ioutil"
	"net/http"
)

//go:generate counterfeiter . CouponService
type CouponService interface {
	CreateCoupon(couponInstance coupon.Coupon) error
	UpdateCoupon(couponInstance coupon.Coupon) error
}

//go:generate counterfeiter . CouponSerializer
type CouponSerializer interface {
	Deserialize(bodyBytes []byte) (coupon.Coupon, error)
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

	couponInstance, err := h.Serializer.Deserialize(bodyBytes)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// TODO MS - add validation step here on couponInstance to assert that all fields are provided
	// otherwise return an error

	err = h.CouponService.CreateCoupon(couponInstance)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h CouponHandler) handlePatch(w http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	couponInstance, err := h.Serializer.Deserialize(bodyBytes)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.CouponService.UpdateCoupon(couponInstance)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}