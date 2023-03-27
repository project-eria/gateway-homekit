package lib

/*
var servicesTempSensor = map[string]*service.TemperatureSensor{}

func initThermometer(info accessory.Info, addrXAAL string) *accessory.Accessory {
	acc := accessory.NewTemperatureSensor(info, 0, -5, 40, 1)

	servicesTempSensor[addrXAAL] = acc.TempSensor
	return acc.Accessory
}

func thermometerSendHomekit(addrXAAL string, values map[string]interface{}) {
	temperature, ok := values["temperature"]
	if ok {
		value := temperature.(float64)
		servicesTempSensor[addrXAAL].CurrentTemperature.SetValue(value)
	} else {
		log.Error().Str("module", "main:thermometer").Error("value 'temperature' not found")
	}
}
*/
