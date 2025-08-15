package utils

func GetCollToken(id int) string {
	if id == 0 {
		return "DAI"
	}

	if id == 1 {
		return "WETH"
	}

	if id == 2 {
		return "USDC"
	}
	return "ERROR PLEASE CHECK GET COLL TOKEN"
}
