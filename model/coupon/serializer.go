package coupon

import (
	"encoding/json"
	"io"
)

type Serializer struct {}

func (s Serializer) Deserialize(body io.ReadCloser) Coupon {
	var coupon Coupon

	decoder := json.NewDecoder(body)
	decoder.Decode(&coupon)

	return coupon
}
