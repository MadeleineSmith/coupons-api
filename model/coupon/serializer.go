package coupon

import (
	"bytes"
	"github.com/google/jsonapi"
)

type Serializer struct {}

func (s Serializer) Deserialize(body []byte) (Coupon, error) {
	coupon := new(Coupon)

	err := jsonapi.UnmarshalPayload(bytes.NewReader(body), coupon)
	if err != nil {
		return Coupon{}, err
	}

	return *coupon, nil
}