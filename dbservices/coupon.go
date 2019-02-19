package dbservices

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/madeleinesmith/coupons/model/coupon"
)

type CouponService struct {
	DB *sql.DB
}

func (s CouponService) CreateCoupon(couponInstance coupon.Coupon) (*coupon.Coupon, error) {
	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("coupons").
		Columns("name", "brand", "value").
		Values(*couponInstance.Name, *couponInstance.Brand, *couponInstance.Value).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return nil, err
	}

	err = s.DB.QueryRow(query, args...).Scan(&couponInstance.ID)
	if err != nil {
		return nil, err
	}

	return &couponInstance, nil
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

func (s CouponService) GetCouponByFilter(filterName string, value string) (*coupon.Coupon, error) {
	sqlString, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("id", "name", "brand", "value").
		From("coupons").
		Where(squirrel.Eq{filterName: value}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := s.DB.QueryRow(sqlString, args...)

	var couponInstance coupon.Coupon
	err = row.Scan(&couponInstance.ID, &couponInstance.Name, &couponInstance.Brand, &couponInstance.Value)
	if err != nil {
		return nil, err
	}

	return &couponInstance, nil
}