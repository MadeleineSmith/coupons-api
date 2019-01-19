package dbservices

import (
	"database/sql"
	"github.com/coupons/model"
)

type CouponService struct {
	DB *sql.DB
}

func (s CouponService) CreateCoupon(coupon model.Coupon) {
	s.DB.Exec("INSERT INTO coupons(name, brand, value) VALUES ($1, $2, $3)", coupon.Name, coupon.Brand, coupon.Value)
}