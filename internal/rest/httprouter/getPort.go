package httprouter

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

/*
Чтение порта из конфига "../../configs/portConfig.json"

:return: порт API или ошибка, если конфиг не удалось прочитать
*/
func getPort() (string, error) {
	// Структура конфигурации порта
	type portConfig struct {
		Port string `json:"port"` // номер порта
	}
	configFilePath := "../../configs/portConfig.json"
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return "", errors.New("fail to read port config " + configFilePath)
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	var portStruct portConfig
	json.Unmarshal(byteValue, &portStruct)

	return portStruct.Port, nil
}
