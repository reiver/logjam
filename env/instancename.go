package env

import (
	"os"
)

var InstanceName string = instanceName()

func instanceName() string {
	value := os.Getenv("INSTANCE_NAME")
	if "" == value  {
		value = "GreatApe"
	}

	return value
}
