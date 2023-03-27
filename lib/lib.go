package lib

import (
	"errors"

	"github.com/brutella/hap/accessory"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

type capability interface {
	setCapability()
	// setConn(*consumer.ConsumedThing)
}

// type thingAccessory struct {
// 	name string
// 	conn *consumer.ThingConnection
// }

// func (t *thingAccessory) setConn(conn *consumer.ThingConnection) {
// 	t.conn = conn
// }

//	a := accessory.NewSwitch(accessory.Info{
//		Name: "Lamp",
//	})
func NewAccessory(t *consumer.ConsumedThing, deviceType string, deviceProperty string) (*accessory.A, error) {
	details := t.GetThingDescription()
	info := accessory.Info{
		Name:         details.Description,
		Manufacturer: "ERIA",
		SerialNumber: details.ID,
	}

	var (
		a *accessory.A
	)
	zlog.Info().Str("label", details.Description).Str("type", deviceType).Str("property", deviceProperty).Msg("[main] Init Accessory")

	switch deviceType {
	case "ShutterBasic":
		_, a = newShutterBasic(info, t)
	case "ShutterPosition":
		_, a = newShutterPosition(info, t)
	case "LampBasic":
		_, a = newLampBasic(info, t)
	case "LampDimmer":
		_, a = newLampDimmer(info, t)
	case "HomeContext":
		_, a = newHomeContext(info, t, deviceProperty)
	/*
		//	case "thermometer.basic":
			return initThermometer(info, conf.Href), nil
	*/
	default:
		return nil, errors.New(deviceType + " type methods hasn't been implemented yet")
	}

	return a, nil
}
