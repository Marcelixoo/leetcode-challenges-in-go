package coinchange

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCase struct {
	coins  []int
	target int
	output int
}

func (c *TestCase) asString() string {
	coinsAsSliceOfString := []string{}
	for i := range c.coins {
		coinAsNumber := c.coins[i]
		coinAsString := strconv.Itoa(coinAsNumber)
		coinsAsSliceOfString = append(coinsAsSliceOfString, coinAsString)
	}

	coinsAsString := strings.Join(coinsAsSliceOfString, ", ")

	return fmt.Sprintf("coins: %v, target: %d, expected_output: %v", coinsAsString, c.target, c.output)
}

func TestCases(t *testing.T) {
	exampleTestCases := []TestCase{
		{
			coins:  []int{},
			target: 0,
			output: 0,
		},
		{
			coins:  []int{6, 5, 1},
			target: 14,
			output: 4,
		},
		{
			coins:  []int{14, 13, 11},
			target: 1,
			output: -1,
		},
		{
			coins:  []int{4, 3, 1},
			target: 6,
			output: 2,
		},
		{
			coins:  []int{25, 10, 5},
			target: 30,
			output: 2,
		},
	}

	var testCase TestCase
	for i := range exampleTestCases {
		testCase = exampleTestCases[i]
		result := MinNumberOfCoinsToReachAmount(testCase.coins, testCase.target)
		require.Equal(t, testCase.output, result, testCase.asString())
	}
}
