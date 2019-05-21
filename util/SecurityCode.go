package util

import (
	"fmt"
	"math/rand"
	"time"
)

func SixRandomNumberByPhone() string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", seed.Int31n(999999))
}

