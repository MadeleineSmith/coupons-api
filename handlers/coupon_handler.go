package handlers

import (
	"fmt"
	"net/http"
)

type CouponHandler struct {}

func (h CouponHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ENDPOINT HIT")
}
