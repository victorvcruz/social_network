package utils

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func GetStringEnvOrElse(envName string, defaultValue string) string {
	value, found := os.LookupEnv(envName)
	if !found {
		value = defaultValue
	}
	return value
}

func GetIntEnvOrElse(envName string, defaultValue int) (value int) {
	valueStr, found := os.LookupEnv(envName)
	if !found {
		value = defaultValue
	} else {
		value, _ = strconv.Atoi(valueStr)
	}
	return value
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func ReadBodyAndReturnMapBody(body io.ReadCloser) (map[string]interface{}, error) {

	bodyByte, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	var mapBody map[string]interface{}

	if err := json.Unmarshal(bodyByte, &mapBody); err != nil {
		return nil, err
	}

	return mapBody, nil
}

func StringNullable(str interface{}) string {
	if str == nil {
		return ""
	}
	return str.(string)
}
