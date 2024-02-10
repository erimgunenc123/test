package environment

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var allowMigration bool
var workEnv string

func IsTestEnvironment() bool {
	return workEnv == TestEnv
}

func AllowMigrations() bool {
	return allowMigration
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
		case "-migrate":
			if strings.ToLower(arg[1]) == "true" {
				allowMigration = true
			}
		default:
			log.Printf("Unexpected arg: %s", arg_)
		}
	}
}
