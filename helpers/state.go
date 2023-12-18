package helpers

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetViperValueEnsureSet(key string) string {
	if !viper.IsSet(key) {
		panic(fmt.Errorf("'%s' needs to be set by running 'configure'", key))
	}

	return viper.GetString(key)
}
