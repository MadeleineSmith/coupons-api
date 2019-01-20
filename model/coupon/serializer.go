package coupon

import (
	"encoding/json"
)

type Serializer struct {}

func (s Serializer) Deserialize(body []byte) (Coupon, error) {
	var coupon Coupon

	json.Unmarshal(body, &coupon)

	return coupon, nil
}