package main

import (
	"github.com/brutella/hc/accessory"
	logger "github.com/project-eria/eria-logger"
)

var windowCoveringList = map[string]*accessory.WindowCovering{}

func initShutter(info accessory.Info, addrXAAL string) *accessory.Accessory {
	acc := accessory.NewWindowCovering(info, 100, 0, 100, 0)
	acc.WindowCovering.CurrentPosition.SetValue(100) // TODO set with current REAL position
	acc.WindowCovering.TargetPosition.SetValue(100)

	acc.WindowCovering.TargetPosition.OnValueRemoteUpdate(func(position int) {
		shutterSendXAALUpDown(addrXAAL, position)
	})
	windowCoveringList[addrXAAL] = acc
	return acc.Accessory
}

func initShutterPosition(info accessory.Info, addrXAAL string) *accessory.Accessory {
	acc := accessory.NewWindowCovering(info, 100, 0, 100, 10)
	acc.WindowCovering.CurrentPosition.SetValue(100) // TODO set with current REAL position
	acc.WindowCovering.TargetPosition.SetValue(100)

	acc.WindowCovering.TargetPosition.OnValueRemoteUpdate(func(position int) {
		shutterSendXAALPosition(addrXAAL, position)
	})
	windowCoveringList[addrXAAL] = acc
	return acc.Accessory
}

func shutterSendXAALUpDown(addrXAAL string, value int) {
	if value > 0 {
		sendXAAL(addrXAAL, "up", nil)
		windowCoveringList[addrXAAL].WindowCovering.TargetPosition.SetValue(100)
	} else {
		sendXAAL(addrXAAL, "down", nil)
		windowCoveringList[addrXAAL].WindowCovering.TargetPosition.SetValue(0)
	}
}

func shutterSendXAALPosition(addrXAAL string, value int) {
	if value == 100 {
		sendXAAL(addrXAAL, "up", nil)
	} else if value == 0 {
		sendXAAL(addrXAAL, "down", nil)
	} else {
		var body = map[string]interface{}{"target": value}
		sendXAAL(addrXAAL, "position", body)
	}
	windowCoveringList[addrXAAL].WindowCovering.TargetPosition.SetValue(value)
}

func shutterSendHomekitUpDown(addrXAAL string, values map[string]interface{}) {
	action, ok := values["action"]
	if ok {
		if action == "up" {
			windowCoveringList[addrXAAL].WindowCovering.CurrentPosition.SetValue(100)
			windowCoveringList[addrXAAL].WindowCovering.TargetPosition.SetValue(100)
		} else {
			windowCoveringList[addrXAAL].WindowCovering.CurrentPosition.SetValue(0)
			windowCoveringList[addrXAAL].WindowCovering.TargetPosition.SetValue(0)
		}
	} else {
		logger.Module("main:shutter").Error("Value 'action' not found")
	}
}

func shutterSendHomekitPosition(addrXAAL string, values map[string]interface{}) {
	_, ok := values["action"]
	if ok {
		logger.Module("main:shutter").Debug("Received 'action' notification (ignored for shutter.position")
		return // Ignore action and wait for final position
	}
	value, ok := values["position"]
	if ok {
		position := int(value.(float64))
		logger.Module("main:shutter").WithField("position", position).Debug("Received 'position' notification")
		windowCoveringList[addrXAAL].WindowCovering.CurrentPosition.SetValue(position)
		windowCoveringList[addrXAAL].WindowCovering.TargetPosition.SetValue(position)
	} else {
		logger.Module("main:shutter").Error("Value 'position' not found")
	}
}
