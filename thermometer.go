package main

import (
	"github.com/brutella/hc/accessory"
	logger "github.com/project-eria/eria-logger"
)

var temperatureSensorList = map[string]*accessory.Thermometer{}

func initThermometer(info accessory.Info, addrXAAL string) *accessory.Accessory {
	acc := accessory.NewTemperatureSensor(info, 0, -5, 40, 1)

	temperatureSensorList[addrXAAL] = acc
	return acc.Accessory
}

func thermometerSendHomekit(addrXAAL string, values map[string]interface{}) {
	temperature, ok := values["temperature"]
	if ok {
		value := temperature.(float64)
		temperatureSensorList[addrXAAL].TempSensor.CurrentTemperature.SetValue(value)
	} else {
		logger.Module("main:thermometer").Error("value 'temperature' not found")
	}
}
