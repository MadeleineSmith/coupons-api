package dbservices

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
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

func (s CouponService) UpdateCoupon(coupon coupon.Coupon) error {
	updateStatement := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Update("coupons").
		Where(squirrel.Eq{"id": coupon.ID})

	if coupon.Name != nil {
		updateStatement = updateStatement.Set("name", &coupon.Name)
	}

	if coupon.Brand != nil {
		updateStatement = updateStatement.Set("brand", &coupon.Brand)
	}

	if coupon.Value != nil {
		updateStatement = updateStatement.Set("value", &coupon.Value)
	}

	dbQuery, args, err := updateStatement.ToSql()
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(dbQuery, args...)
	if err != nil {
		return err
	}

	return nil
}