package validators

import (
	"errors"
	"github.com/madeleinesmith/coupons/model/coupon"
	"strings"
)

type CouponValidator struct{}

func (v CouponValidator) isEmptyField(fieldValue string) bool {
	return (len(strings.Trim(fieldValue, " \n\t"))) < 1
}

func (v CouponValidator) Validate(coupon coupon.Coupon) error {
	if coupon.Name == nil || v.isEmptyField(*coupon.Name) {
		return errors.New("name field is required")
	}

	if coupon.Brand == nil || v.isEmptyField(*coupon.Brand) {
		return errors.New("brand field is required")
	}

	if coupon.Value == nil {
		return errors.New("value field is required")
	}

	return nil
}
