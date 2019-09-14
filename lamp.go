package main

import (
	"github.com/brutella/hc/accessory"
	logger "github.com/project-eria/eria-logger"
)

var lightbulbList = map[string]*accessory.Lightbulb{}

func initLamp(info accessory.Info, addrXAAL string) *accessory.Accessory {
	acc := accessory.NewLightbulb(info)

	acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		lampSendXAALOnOff(addrXAAL, on)
	})
	lightbulbList[addrXAAL] = acc
	return acc.Accessory
}

func initLampDimmer(info accessory.Info, addrXAAL string) *accessory.Accessory {
	acc := accessory.NewLightbulb(info)

	acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		lampSendXAALOnOff(addrXAAL, on)
	})

	acc.Lightbulb.Brightness.SetValue(0) // TODO set with current REAL value
	acc.Lightbulb.Brightness.OnValueRemoteUpdate(func(value int) {
		var body = map[string]interface{}{"target": value}
		sendXAAL(addrXAAL, "dim", body)
	})

	lightbulbList[addrXAAL] = acc
	return acc.Accessory
}

func lampSendHomekitOnOff(addrXAAL string, values map[string]interface{}) {
	light, ok := values["light"]
	if ok {
		value := (light == "on")
		lightbulbList[addrXAAL].Lightbulb.On.SetValue(value)
	} else {
		logger.Module("main:lamp").Error("value 'light' not found")
	}
}

func lampSendHomekitDim(addrXAAL string, values map[string]interface{}) {
	light, ok := values["light"]
	if ok {
		value := (light == "on")
		lightbulbList[addrXAAL].Lightbulb.On.SetValue(value)
		return
	}
	value, ok := values["dimmer"]
	if ok {
		dimmer := int(value.(float64))
		logger.Module("main:lamp").WithField("dimmer", dimmer).Debug("Received 'dimmer' notification")
		lightbulbList[addrXAAL].Lightbulb.Brightness.SetValue(dimmer)
	} else {
		logger.Module("main:lamp").Error("Value 'position' not found")
	}
}

func lampSendXAALOnOff(addrXAAL string, value bool) {
	if value {
		sendXAAL(addrXAAL, "on", nil)
	} else {
		sendXAAL(addrXAAL, "off", nil)
	}
}
