package environment

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var workEnv string

func IsTestEnvironment() bool {
	return workEnv == TestEnv
}

func ParseArgs() {
	workEnv = TestEnv
	for _, arg_ := range os.Args[1:] {
		arg := strings.Split(arg_, "=")
		switch arg[0] {
		case "-env":
			if arg[1] == "prod" {
				workEnv = ProdEnv
			} else if arg[1] == "test" {
				workEnv = TestEnv
			} else {
				panic(fmt.Sprintf("invalid arg: %s", arg_))
			}
		default:
			log.Printf("Unexpected arg: %s", arg_)
		}
	}
}
