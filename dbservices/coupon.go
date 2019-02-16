package dbservices

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/madeleinesmith/coupons/model/coupon"
)

type CouponService struct {
	DB *sql.DB
}

func (s CouponService) CreateCoupon(couponInstance coupon.Coupon) (coupon.Coupon, error) {
	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("coupons").
		Columns("name", "brand", "value").
		Values(*couponInstance.Name, *couponInstance.Brand, *couponInstance.Value).
		Suffix("RETURNING id").
		ToSql()

	// able to test this error case?
	// right to return empty coupon here?
	// and below also?
	// + testing strategy -> mocked db/ real db # of tests
	if err != nil {
		return coupon.Coupon{}, err
	}

	err = s.DB.QueryRow(query, args...).Scan(&couponInstance.ID)
	if err != nil {
		return coupon.Coupon{}, err
	}

	return couponInstance, nil
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

func (s CouponService) GetCoupons() ([]*coupon.Coupon, error) {
	selectStatement := squirrel.StatementBuilder.
		Select("id, name, brand, value").
		From("coupons")

	dbQuery, _, err := selectStatement.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(dbQuery)
	if err != nil {
		return nil, err
	}

	var couponSlice []*coupon.Coupon

	for rows.Next() {
		couponInstance := new(coupon.Coupon)

		err := rows.Scan(&couponInstance.ID, &couponInstance.Name, &couponInstance.Brand, &couponInstance.Value)

		if err != nil {
			return nil, err
		}

		couponSlice = append(couponSlice, couponInstance)
	}

	return couponSlice, nil
}