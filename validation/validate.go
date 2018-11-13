// Package validation handles the validation of user-provided OEM and CD keys.
package validation

import (
	"fmt"
	"strconv"
)

func checkdigitCheck(c string) bool {
	// Check digit cannot be 0 or >= 8.
	if c[len(c)-1:] == "0" || c[len(c)-1:] == "8" || c[len(c)-1:] == "9" {
		return false
	}
	return true
}

func digitsum(num int64) int64 {
	var s int64
	for num != 0 {
		digit := num % 10
		s += digit
		num /= 10
	}
	return s
}

func validateCDKey(key string) error {
	site, err := strconv.ParseInt(key[0:3], 10, 0)
	if err != nil {
		return fmt.Errorf("the site number isn't a number")
	}
	main, err := strconv.ParseInt(key[4:11], 10, 0)
	if err != nil {
		return fmt.Errorf("the second segment isn't a number")
	}

	invalidSites := map[int64]int{333: 333, 444: 444, 555: 555, 666: 666, 777: 777, 888: 888, 999: 999}
	_, invalid := invalidSites[site]
	if invalid {
		fmt.Println("The site number is invalid: site number cannot be 333, 444, 555, 666, 777, 888, or 999.")
	}

	c := strconv.Itoa(int(main))
	// We must check the check digit.
	if !checkdigitCheck(c) {
		fmt.Println("The second segment of the key is invalid: the last digit cannot be 0 or >= 8.")
	}

	// Split the second segment to individual digits for the division check.
	sum := digitsum(main)
	if sum%7 != 0 {
		fmt.Printf("The second segment of the key is invalid: the digit sum (%d) must be divisible by 7.", sum)
	}
	return nil
}

func validateOEM(key string) error {
	_, err := strconv.ParseInt(key[0:5], 10, 0)
	if err != nil {
		return fmt.Errorf("the first segment is not a number")
	}
	th, err := strconv.ParseInt(key[10:17], 10, 0)
	if err != nil {
		return fmt.Errorf("the third segment is not a number")
	}
	julian, err := strconv.ParseInt(key[0:3], 10, 0)
	if julian == 0 || julian > 366 {
		fmt.Println("The date is invalid: date has to be 001-366.")
	}
	year := key[3:5]
	validYears := map[string]string{"95": "95", "96": "96", "97": "97", "98": "98", "99": "99", "00": "00", "01": "01", "02": "02", "03": "03"}
	_, valid := validYears[year]
	if !valid {
		fmt.Println("The year is invalid: cannot be less than 95 or above 03")
	}

	third := key[10:17]
	if string(third[0]) != "0" {
		fmt.Println("The third segment is invalid: must begin with a 0.")
	}
	c := strconv.Itoa(int(th))
	// We must check the check digit.
	if !checkdigitCheck(c) {
		fmt.Println("The third segment of the key is invalid: the last digit cannot be 0 or >= 8.")
	}
	// Split the third segment to individual digits for the division check.
	sum := digitsum(th)
	if sum%7 != 0 {
		fmt.Printf("The third segment of the key is invalid: the digit sum (%d) must be divisible by 7.", sum)
	}
	return nil
}

// ValidateKey validates the provided OEM or CD key.
func ValidateKey(k string) {
	// Make sure the provided key has a chance of being valid.
	switch {
	case len(k) == 11 && k[3:4] == "-":
		fmt.Printf("%s could be a valid CD key. The key is valid if you get no output.\n", k)
		if err := validateCDKey(k); err != nil {
			fmt.Println("Unable to validate key:", err)
		}
	case len(k) == 23 && k[5:6] == "-" && k[9:10] == "-" && k[17:18] == "-" && len(k[18:]) == 5:
		fmt.Printf("%s could be a valid OEM key. The key is valid if you get no output.\n", k)
		if err := validateOEM(k); err != nil {
			fmt.Println("Unable to validate key:", err)
		}
	default:
		fmt.Printf("%s doesn't even resemble a valid key.\n", k)
	}

}
