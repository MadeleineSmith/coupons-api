package handlers

import (
	"errors"
	"github.com/lib/pq"
	"github.com/madeleinesmith/coupons/model/coupon"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Filters struct {
	Name  *string
	Value *int
	Brand *string
}

//go:generate counterfeiter . CouponService
type CouponService interface {
	CreateCoupon(couponInstance coupon.Coupon) (*coupon.Coupon, error)
	UpdateCoupon(couponInstance coupon.Coupon) error
	GetCoupons(filters Filters) ([]*coupon.Coupon, error)
	GetCouponById(couponId string) (*coupon.Coupon, error)
}

//go:generate counterfeiter . CouponSerializer
type CouponSerializer interface {
	DeserializeCoupon(bodyBytes []byte) (coupon.Coupon, error)
	SerializeCoupon(coupon *coupon.Coupon) ([]byte, error)
	SerializeCoupons([]*coupon.Coupon) ([]byte, error)
}

//go:generate counterfeiter . CouponValidator
type CouponValidator interface {
	Validate(coupon coupon.Coupon) error
}

type CouponHandler struct {
	Serializer      CouponSerializer
	CouponService   CouponService
	CouponValidator CouponValidator
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

	err = h.CouponValidator.Validate(couponInstance)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// consider sanitizing the coupon i.e. removing whitespace from fields before inserting into the db
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
	var coupons []*coupon.Coupon
	var err error

	var filters Filters

	for queryParamsKey, queryParamsValue := range req.URL.Query() {
		if queryParamsKey == "brand" {
			brand := queryParamsValue[0]
			filters.Brand = &brand

		} else if queryParamsKey == "value" {
			stringValue := queryParamsValue[0]

			intValue, err := strconv.Atoi(stringValue)
			if err != nil {
				handleError(w, err, http.StatusBadRequest)
				return
			}

			filters.Value = &intValue

		} else if queryParamsKey == "name" {
			name := queryParamsValue[0]
			filters.Name = &name
		}
	}

	coupons, err = h.CouponService.GetCoupons(filters)

	if err != nil {
		code := http.StatusInternalServerError

		// 42703 is an undefined_column error
		pqError, ok := err.(*pq.Error)
		if ok {
			if pqError.Code == "42703" {
				code = http.StatusBadRequest
			}
		}

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
