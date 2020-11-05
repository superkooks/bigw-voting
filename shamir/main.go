package shamir

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	fieldSize = int(math.Pow(2, 31) - 1)
)

// Invert number in field n using the extended euclidean algorithm
func extendedEuclideanInverse(a, n int) int {
	t, newT := 0, 1
	r, newR := n, int(math.Abs(float64(a)))

	for newR != 0 {
		quotient := r / newR
		t, newT = newT, t-quotient*newT
		r, newR = newR, r-quotient*newR
	}

	if r > 1 {
		panic("a is not invertible")
	}

	if t < 0 {
		t = t + n
	}

	return t
}

func divideInField(a, b int) int {
	inverseB := extendedEuclideanInverse(b, fieldSize)

	if b < 0 {
		inverseB *= -1
	}
	return ((a % fieldSize) * ((inverseB % fieldSize) % fieldSize)) % fieldSize
}

func lagrangeInterpolate(x int, xValues []int, yValues []int) (int, error) {
	if len(xValues) == 0 || len(xValues) != len(yValues) {
		return 0, fmt.Errorf("x and y cannot be empty and must be equal lengths")
	}

	var sum int
	for j := 0; j < len(xValues); j++ {
		var numerators []int
		var denominators []int
		for m := 0; m < len(xValues); m++ {
			if m != j {
				numerators = append(numerators, (x-xValues[m])%fieldSize)
				denominators = append(denominators, (xValues[j]-xValues[m])%fieldSize)
			}
		}

		numerator := 1
		for _, val := range numerators {
			numerator *= val
		}

		denominator := 1
		for _, val := range denominators {
			denominator *= val
		}

		sum = (sum + int(divideInField(numerator, denominator))*yValues[j]) % fieldSize
	}

	return (sum + fieldSize) % fieldSize, nil
}

func reconstructSecret(points [][2]int) (int, error) {
	if len(points) < 2 {
		return 0, fmt.Errorf("at least two shares are required")
	}

	var xValues []int
	var yValues []int

	for _, pair := range points {
		xValues = append(xValues, pair[0])
		yValues = append(yValues, pair[1])
	}

	s, _ := lagrangeInterpolate(0, xValues, yValues)
	return s % fieldSize, nil
}

func constructPoints(secret int, shares int, sufficientShares int) [][2]int {
	rand.Seed(time.Now().UnixNano())
	randomNumbers := make([]int, sufficientShares-1)
	for i := 0; i < sufficientShares-1; i++ {
		randomNumbers[i] = rand.Intn(int(fieldSize))
	}

	out := make([][2]int, shares)
	for index := range out {
		point := secret
		for xPower := 1; xPower < sufficientShares; xPower++ {
			x := int(math.Pow(float64(index+1), float64(xPower)))
			point = (point + randomNumbers[xPower-1]*x) % int(fieldSize)
		}

		out[index] = [2]int{index + 1, point}
	}

	return out
}
