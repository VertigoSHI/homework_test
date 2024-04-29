package util

import "fmt"

func IntToHexString(num int) string {
	return fmt.Sprintf("0x%x", num)
}
