package shamir

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	// FieldSize is the size of the integers mod p field
	FieldSize = int(math.Pow(2, 31) - 1)
)

// Invert number in field n using the extended euclidean algorithm
func extendedEuclideanInverse(a, n int) (int, error) {
	t, newT := 0, 1
	r, newR := n, int(math.Abs(float64(a)))

	for newR != 0 {
		quotient := r / newR
		t, newT = newT, t-quotient*newT
		r, newR = newR, r-quotient*newR
	}

	if r > 1 {
		return 0, errors.New("a is not invertible")
	}

	if t < 0 {
		t = t + n
	}

	return t, nil
}

func divideInField(a, b int) (int, error) {
	inverseB, err := extendedEuclideanInverse(b, FieldSize)
	if err != nil {
		return 0, err
	}

	if b < 0 {
		inverseB *= -1
	}

	return ((a % FieldSize) * ((inverseB % FieldSize) % FieldSize)) % FieldSize, nil
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
				numerators = append(numerators, (x-xValues[m])%FieldSize)
				denominators = append(denominators, (xValues[j]-xValues[m])%FieldSize)
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

		div, err := divideInField(numerator, denominator)
		if err != nil {
			return 0, nil
		}

		sum = (sum + int(div)*yValues[j]) % FieldSize
	}

	return (sum + FieldSize) % FieldSize, nil
}

// ReconstructSecret reconstructs the secret from a given set of points
func ReconstructSecret(points [][2]int) (int, error) {
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
	return s % FieldSize, nil
}

// ConstructPoints constructs a new Shamir polynomial and distributes the points
func ConstructPoints(secret int, shares int, sufficientShares int) [][2]int {
	rand.Seed(time.Now().UnixNano())
	randomNumbers := make([]int, sufficientShares-1)
	for i := 0; i < sufficientShares-1; i++ {
		randomNumbers[i] = rand.Intn(int(FieldSize))
	}

	out := make([][2]int, shares)
	for index := range out {
		point := secret
		for xPower := 1; xPower < sufficientShares; xPower++ {
			x := int(math.Pow(float64(index+1), float64(xPower)))
			point = (point + randomNumbers[xPower-1]*x) % int(FieldSize)
		}

		out[index] = [2]int{index + 1, point}
	}

	return out
}
