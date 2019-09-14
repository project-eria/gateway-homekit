package main

import (
	"fmt"

	"github.com/brutella/hc/accessory"
	logger "github.com/project-eria/eria-logger"
	"github.com/project-eria/xaal-go"
	"github.com/project-eria/xaal-go/message"
)

func initAccessory(conf *configDevice) (*accessory.Accessory, error) {
	info := accessory.Info{
		Name:             conf.Name,
		Manufacturer:     "ERIA",
		SerialNumber:     conf.XaalAddr,
		FirmwareRevision: Version,
	}

	switch conf.Type {
	case "lamp.basic":
		return initLamp(info, conf.XaalAddr), nil
	case "lamp.dimmer":
		return initLampDimmer(info, conf.XaalAddr), nil
	case "shutter.basic":
		return initShutter(info, conf.XaalAddr), nil
	case "shutter.position":
		return initShutterPosition(info, conf.XaalAddr), nil
	case "thermometer.basic":
		return initThermometer(info, conf.XaalAddr), nil
	default:
	}
	return nil, fmt.Errorf("%s type methods hasn't been implemented yet", conf.Type)
}

func sendXAAL(address string, action string, args map[string]interface{}) {
	addresses := []string{address}
	go xaal.SendRequest(_gw, addresses, action, args)
}

func updateFromXAAL(msg *message.Message) {
	// TODO send Xaal notifications to Homekit
	_, ok := _devicesByXAAL[msg.Header.Source]
	if ok && msg.IsAttributesChange() {
		switch msg.Header.DevType {
		case "lamp.basic":
			lampSendHomekitOnOff(msg.Header.Source, msg.Body)
			break
		case "lamp.dimmer":
			lampSendHomekitDim(msg.Header.Source, msg.Body)
			break
		case "shutter.basic":
			shutterSendHomekitUpDown(msg.Header.Source, msg.Body)
			break
		case "shutter.position":
			shutterSendHomekitPosition(msg.Header.Source, msg.Body)
			break
		case "thermometer.basic":
			thermometerSendHomekit(msg.Header.Source, msg.Body)
			break
		default:
			// TODO return fmt.Errorf("%s type methods hasn't been implemented yet", typeXAAL)
		}
		logger.Module("main").WithField("msg", msg).Trace()
	}
}

func addrInDevices(addr string) *configDevice {
	for i := range config.Devices {
		confDev := &config.Devices[i]

		if confDev.XaalAddr == addr {
			return confDev
		}
	}
	return nil
}
