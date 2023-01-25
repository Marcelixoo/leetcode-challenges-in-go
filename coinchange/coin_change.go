package coinchange

// Please discard amounts greater than this number
const A_HUGE_NUMBER = 100000

func Explain() string {
	return "322. Coin Change: https://leetcode.com/problems/coin-change"
}

func MinNumberOfCoinsToReachAmount(coins []int, amount int) int {
	/*
		In a nutschel, the idea behind the solution
		is to store previous values in order to find
		the next ones.

		The formula seem to be a result of empirical analysis, which makes it hard to describe
		into words. So let's look at an example!

		For a set of coins [4, 3, 1] we have the
		following values:

		|0|1|2|3|4|5|6|...|N
		|0|1|2|1|1|2|2|...|MIN(v[i], v[i - coin] + 1)
	*/

	if amount == 0 {
		return 0
	}

	mapOfValuesPerAmount := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		mapOfValuesPerAmount[i] = A_HUGE_NUMBER

		for _, coin := range coins {
			if coin <= i {
				mapOfValuesPerAmount[i] = min(mapOfValuesPerAmount[i], mapOfValuesPerAmount[i-coin]+1)
			}
		}
	}

	valueForTargetAmount := mapOfValuesPerAmount[amount]
	if valueForTargetAmount == A_HUGE_NUMBER {
		/*
			Meaning: if we could not find a plausible number of coins to achieve the given `amount`
		*/
		return -1
	}
	return valueForTargetAmount
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MinNumberOfCoinsToReachAmountV0(coins []int, amount int) int {
	/*
		My naive first solution simply returned
		the first combination possible of numbers,
		taking advantage of the ordering imposed
		to the list of coins.

		It obviously failed to accomodate all
		test cases but still helped me wrap my
		head around the problem :)
	*/

	if amount == 0 {
		return 0
	}
	return recurse(coins, amount, 0)
}

func recurse(coins []int, amount int, count int) int {
	if amount == 0 {
		return count
	}

	var newAmount int
	var coin int

	for i := range coins {
		coin = coins[i]
		newAmount = amount - coin

		if newAmount >= 0 {
			return recurse(coins, newAmount, count+1)
		}
	}
	return -1
}
