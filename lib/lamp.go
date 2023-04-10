package lib

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

type lampBasic struct {
	*consumer.ConsumedThing
	accessory *accessory.Lightbulb
}

func newLampBasic(info accessory.Info, t *consumer.ConsumedThing) (*lampBasic, *accessory.A) {
	acc := accessory.NewLightbulb(info)
	/* 'on' property */
	setOn(acc.Lightbulb.On, info, t)
	return &lampBasic{accessory: acc, ConsumedThing: t}, acc.A
}

func (d lampBasic) setCapability() {}

type DimmingLightbulb struct {
	*accessory.A
	Lightbulb *serviceDimmingLightbulb
}

func NewDimmingLightbulb(info accessory.Info) *DimmingLightbulb {
	a := DimmingLightbulb{}
	a.A = accessory.New(info, accessory.TypeLightbulb)
	a.Lightbulb = NewServiceDimmingLightbulb()
	a.AddS(a.Lightbulb.S)

	return &a
}

type serviceDimmingLightbulb struct {
	*service.S
	On         *characteristic.On
	Brightness *characteristic.Brightness
}

func NewServiceDimmingLightbulb() *serviceDimmingLightbulb {
	s := serviceDimmingLightbulb{}
	s.S = service.New(service.TypeLightbulb)

	s.On = characteristic.NewOn()
	s.AddC(s.On.C)

	s.Brightness = characteristic.NewBrightness()
	s.AddC(s.Brightness.C)

	return &s
}

type lampDimmer struct {
	*consumer.ConsumedThing
	accessory *DimmingLightbulb
}

func newLampDimmer(info accessory.Info, t *consumer.ConsumedThing) (*lampDimmer, *accessory.A) {
	acc := NewDimmingLightbulb(info)
	/* 'on' property */
	setOn(acc.Lightbulb.On, info, t)
	/* 'brightness' property */
	setBrightness(acc.Lightbulb.Brightness, info, t)
	return &lampDimmer{accessory: acc, ConsumedThing: t}, acc.A
}

func (d lampDimmer) setCapability() {}

func setOn(accOn *characteristic.On, info accessory.Info, t *consumer.ConsumedThing) {
	data, err := t.ReadProperty("on")
	if err != nil {
		zlog.Error().Err(err).Msg("[main] Can't read lamp state")
	} else {
		state := data.(bool)
		zlog.Trace().Str("name", info.Name).Bool("value", state).Msg("[main] Set lamp initial value")
		accOn.SetValue(state)
	}
	accOn.OnValueRemoteUpdate(func(state bool) {
		zlog.Trace().Str("name", info.Name).Bool("on", state).Msg("[main] Received Lamp update from Homekit")
		t.InvokeAction("toggle", nil)
	})
	t.ObserveProperty("on", func(value interface{}, err error) {
		on := value.(bool)
		zlog.Trace().Str("name", info.Name).Bool("on", on).Msg("[main] Received Lamp update from thing device")
		accOn.SetValue(on)
	})
}

func setBrightness(accBrightness *characteristic.Brightness, info accessory.Info, t *consumer.ConsumedThing) {
	data, err := t.ReadProperty("brightness")
	if err != nil {
		zlog.Error().Err(err).Msg("[main] Can't read lamp brightness")
	} else {
		brightness := int(data.(float64))
		accBrightness.SetValue(brightness)
	}
	accBrightness.OnValueRemoteUpdate(func(brightness int) {
		zlog.Trace().Str("name", info.Name).Int("brightness", brightness).Msg("[main] Received Lamp update from Homekit")
		t.InvokeAction("fade", brightness)
	})
	t.ObserveProperty("brightness", func(value interface{}, err error) {
		brightness := int(value.(float64))
		zlog.Trace().Str("name", info.Name).Int("brightness", brightness).Msg("[main] Received Lamp update from thing device")
		accBrightness.SetValue(brightness)
	})
}
