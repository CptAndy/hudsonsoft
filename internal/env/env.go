package env

import (
	"os"
	"strconv"
)

// gets evn vars
// recives key and a fallback
func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	//Going to recive a value to create an int
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt

}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolValue, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolValue
}
