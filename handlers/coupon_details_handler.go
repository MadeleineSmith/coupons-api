package handlers

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type CouponDetailsHandler struct {
	CouponService CouponService
	Serializer    CouponSerializer
}

func (h CouponDetailsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.handleGet(w, req)
	default:
		err := errors.New(`Method not allowed`)
		handleError(w, err, http.StatusMethodNotAllowed)
	}
}

func (h CouponDetailsHandler) handleGet(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var couponId string
	var ok bool

	if couponId, ok = vars["couponId"]; !ok {
		err := errors.New("couponId URL variable not found")
		handleError(w, err, http.StatusBadRequest)
		return
	}

	couponInstance, err := h.CouponService.GetCouponById(couponId)
	if err != nil {
		code := http.StatusInternalServerError

		if err == sql.ErrNoRows {
			code = http.StatusNotFound
		}

		handleError(w, err, code)
		return
	}

	serializedCoupon, err := h.Serializer.SerializeCoupon(couponInstance)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(serializedCoupon)
}
