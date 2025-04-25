package common

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

func ViperGetStringSlice(key string) ([]string, error) {
	// viper.GetStringSlice does not work properly https://github.com/spf13/viper/issues/380#issuecomment-445244338
	var res []string
	err := viper.UnmarshalKey(key, &res)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("error reading %s from env", key))
	} else {
		return res, nil
	}
}
