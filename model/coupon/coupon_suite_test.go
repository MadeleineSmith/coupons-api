package coupon_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCoupon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Coupon Suite")
}
