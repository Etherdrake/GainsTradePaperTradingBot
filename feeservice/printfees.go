package feeservice

import "fmt"

func TestFees() {
	chain := "arb"
	isLong := false
	collat := 1000
	lev := 30
	pairIndex := 3

	feeData, err := GetFees(chain, isLong, collat, lev, pairIndex)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Borrowing Fee per Block: %.6f\n", feeData)
	fmt.Printf("Borrowing Fee per Hour: %.6f\n", feeData.Borrowing.FeePerHour)
	fmt.Printf("Fees: %d\n", feeData.Fees)
	fmt.Printf("Average Spread: %.4f\n", feeData.AvgSpread)
}
