package coupon

type Coupon struct {
	ID string `jsonapi:"primary,coupons"`
	Name string `jsonapi:"attr,name,omitempty"`
	Brand string `jsonapi:"attr,brand,omitempty"`
	Value int `jsonapi:"attr,value,omitempty"`
}