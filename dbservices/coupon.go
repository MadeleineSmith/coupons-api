package dbservices

import (
	"database/sql"
	"github.com/madeleinesmith/coupons/model/coupon"
)

type CouponService struct {
	DB *sql.DB
}

func (s CouponService) CreateCoupon(coupon coupon.Coupon) error {
	_, err := s.DB.Exec("INSERT INTO coupons (name, brand, value) VALUES ($1, $2, $3)", coupon.Name, coupon.Brand, coupon.Value)
	if err != nil {
		return err
	}

	return nil
}