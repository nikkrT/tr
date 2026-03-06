package paymentService

import (
	"math/rand"
)

func CheckPayment(orderId int) bool {
	return rand.Intn(10) < 7
}
