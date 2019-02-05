package coupon

import (
	"bufio"
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

func (s Serializer) Serialize(coupon Coupon) ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	err := jsonapi.MarshalPayload(writer, &coupon)
	if err != nil {
		return nil, err
	}

	writer.Flush()

	return buffer.Bytes(), nil
}