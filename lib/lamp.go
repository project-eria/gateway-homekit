package lib

import (
	"github.com/brutella/hap/accessory"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

type lampBasic struct {
	accessory *accessory.Lightbulb
	*consumer.ConsumedThing
}

func newLampBasic(info accessory.Info, t *consumer.ConsumedThing) (*lampBasic, *accessory.A) {
	acc := accessory.NewLightbulb(info)
	data, err := t.ReadProperty("on")
	if err != nil {
		zlog.Error().Err(err).Msg("[main] Can't read lamp state")
	} else {
		state := data.(bool)
		acc.Lightbulb.On.SetValue(state)
	}
	acc.Lightbulb.On.OnValueRemoteUpdate(func(state bool) {
		zlog.Trace().Str("name", info.Name).Bool("on", state).Msg("[main] Received update from Homekit")
		t.InvokeAction("toggle", nil)
	})

	t.ObserveProperty("on", func(value interface{}, err error) {
		on := value.(bool)
		zlog.Trace().Str("name", info.Name).Bool("on", on).Msg("[main] Received update from thing device")
		acc.Lightbulb.On.SetValue(on)
	})
	return &lampBasic{accessory: acc, ConsumedThing: t}, acc.A
}

func (d lampBasic) setCapability() {
}

type lampDimmer struct {
	*consumer.ConsumedThing
	accessory *accessory.ColoredLightbulb
}

func newLampDimmer(info accessory.Info, t *consumer.ConsumedThing) (*lampDimmer, *accessory.A) {
	acc := accessory.NewColoredLightbulb(info)
	data, err := t.ReadProperty("on")
	if err != nil {
		zlog.Error().Err(err).Msg("[main] Can't read lamp state")
	} else {
		state := data.(bool)
		acc.Lightbulb.On.SetValue(state)
	}
	acc.Lightbulb.On.OnValueRemoteUpdate(func(state bool) {
		t.InvokeAction("toggle", nil)
	})
	data, err = t.ReadProperty("brightness")
	if err != nil {
		zlog.Error().Err(err).Msg("[main] Can't read lamp brightness")
	} else {
		brightness := int(data.(float64))
		acc.Lightbulb.Brightness.SetValue(brightness)
	}
	acc.Lightbulb.Brightness.OnValueRemoteUpdate(func(brightness int) {
		t.InvokeAction("fade", brightness)
	})
	return &lampDimmer{accessory: acc, ConsumedThing: t}, acc.A
}

func (d lampDimmer) setCapability() {
}

/*
var servicesLightbulb = map[string]*service.Lightbulb{}
var servicesColoredLightbulb = map[string]*service.ColoredLightbulb{}

func initLamp(info accessory.Info, addrXAAL string) *accessory.Accessory {
	acc := accessory.NewLightbulb(info)

	acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		lampSendXAALOnOff(addrXAAL, on)
	})
	servicesLightbulb[addrXAAL] = acc.Lightbulb
	return acc.Accessory
}

func initLampDimmer(info accessory.Info, addrXAAL string) *accessory.Accessory {
	acc := accessory.NewColoredLightbulb(info)

	acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		lampSendXAALOnOff(addrXAAL, on)
	})

	acc.Lightbulb.Brightness.SetValue(0) // TODO set with current REAL value
	acc.Lightbulb.Brightness.OnValueRemoteUpdate(func(value int) {
		var body = map[string]interface{}{"target": value}
		sendXAAL(addrXAAL, "dim", body)
	})

	servicesColoredLightbulb[addrXAAL] = acc.Lightbulb
	return acc.Accessory
}

func lampSendHomekitOnOff(addrXAAL string, values map[string]interface{}) {
	light, ok := values["light"]
	if ok {
		value := (light == "on")
		servicesLightbulb[addrXAAL].On.SetValue(value)
	} else {
		log.Str("module", "main:lamp").Error("value 'light' not found")
	}
}

func lampSendHomekitDim(addrXAAL string, values map[string]interface{}) {
	light, ok := values["light"]
	if ok {
		value := (light == "on")
		servicesColoredLightbulb[addrXAAL].On.SetValue(value)
		return
	}
	value, ok := values["dimmer"]
	if ok {
		dimmer := int(value.(float64))
		log.Str("module", "main:lamp").Str("dimmer", dimmer).Debug("Received 'dimmer' notification")
		servicesColoredLightbulb[addrXAAL].Brightness.SetValue(dimmer)
	} else {
		log.Str("module", "main:lamp").Error("Value 'position' not found")
	}
}

func lampSendXAALOnOff(addrXAAL string, value bool) {
	if value {
		sendXAAL(addrXAAL, "on", nil)
	} else {
		sendXAAL(addrXAAL, "off", nil)
	}
}
*/
