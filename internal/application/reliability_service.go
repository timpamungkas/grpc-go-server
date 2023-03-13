package application

import (
	"fmt"
	"math/rand"
	"time"
)

type ReliabilityService struct {
}

func (r *ReliabilityService) GenerateReliability(minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32) (string, uint32) {
	delay := rand.Intn(int(maxDelaySecond - minDelaySecond + 1 + minDelaySecond))
	delaySecond := time.Duration(delay) * time.Second
	time.Sleep(delaySecond)

	idx := rand.Intn(len(statusCodes))
	str := fmt.Sprintf("The time now is %v, execution delayed for %v seconds",
		time.Now().Format("15:04:05.000"), delay)

	return str, statusCodes[idx]
}
