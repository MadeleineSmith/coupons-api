package coupon

import (
	"encoding/json"
)

type Serializer struct {}

func (s Serializer) Deserialize(body []byte) (Coupon, error) {
	var coupon Coupon

	err := json.Unmarshal(body, &coupon)
	if err != nil {
		return Coupon{}, err
	}

	return coupon, nil
}