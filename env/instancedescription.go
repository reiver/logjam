package env

import (
	"os"
)

var InstanceDescription string = instanceDescription()

func instanceDescription() string {
	value := os.Getenv("INSTANCE_DESCRIPTION")
	if "" == value  {
		value = "A GreatApe instance server — behold the mighty GreatApe"
	}

	return value
}
