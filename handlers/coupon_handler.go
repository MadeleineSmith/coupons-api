package handlers

import (
	"errors"
	"github.com/madeleinesmith/coupons/model/coupon"
	"io/ioutil"
	"net/http"
)

type Filter struct {
	FilterName  string
	FilterValue string
}

//go:generate counterfeiter . CouponService
type CouponService interface {
	CreateCoupon(couponInstance coupon.Coupon) (*coupon.Coupon, error)
	UpdateCoupon(couponInstance coupon.Coupon) error
	GetCoupons(filters ...Filter) ([]*coupon.Coupon, error)
	GetCouponById(couponId string) (*coupon.Coupon, error)
}

//go:generate counterfeiter . CouponSerializer
type CouponSerializer interface {
	DeserializeCoupon(bodyBytes []byte) (coupon.Coupon, error)
	SerializeCoupon(coupon *coupon.Coupon) ([]byte, error)
	SerializeCoupons([]*coupon.Coupon) ([]byte, error)
}

type CouponHandler struct {
	Serializer    CouponSerializer
	CouponService CouponService
}

func (h CouponHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		h.handlePost(w, req)
	case http.MethodPatch:
		h.handlePatch(w, req)
	case http.MethodGet:
		h.handleGet(w, req)
	default:
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
	queryParamsMap := req.URL.Query()
	var filterSlice []Filter

	var coupons []*coupon.Coupon
	var err error

	if len(queryParamsMap) > 0 {
		for key, value := range queryParamsMap {
			var filter Filter
			filter.FilterName = key
			filter.FilterValue = value[0]
			filterSlice = append(filterSlice, filter)
		}
		coupons, err = h.CouponService.GetCoupons(filterSlice...)
	} else if len(queryParamsMap) == 0 {
		coupons, err = h.CouponService.GetCoupons()
	}
	if err != nil {
		code := http.StatusInternalServerError

		if err.Error() == "sql: no rows in result set" {
			code = http.StatusNotFound
		}

		handleError(w, err, code)
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
