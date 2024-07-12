package main

import (
	"fmt"

	regexp "github.com/dlclark/regexp2"
)

func RegEnv() {
	submitEnv := "123"
	RegexUpdate := "z"
	reg, err := regexp.MustCompile(RegexUpdate, regexp.None).FindStringMatch(submitEnv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(reg)
}

func main() {
	RegEnv()
}
