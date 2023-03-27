package lib

import (
	"github.com/brutella/hap/accessory"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

type homeContext struct {
	accessory *accessory.Switch
	*consumer.ConsumedThing
}

func newHomeContext(info accessory.Info, t *consumer.ConsumedThing, p string) (*homeContext, *accessory.A) {
	acc := accessory.NewSwitch(info)
	data, err := t.ReadProperty(p)
	if err != nil {
		zlog.Error().Err(err).Msg("[main] Can't read lamp state")
	} else {
		state := data.(bool)
		acc.Switch.On.SetValue(state)
	}
	acc.Switch.On.OnValueRemoteUpdate(func(state bool) {
		zlog.Trace().Str("name", info.Name).Bool("on", state).Msg("[main] Received update from Homekit")
		t.WriteProperty(p, state)
	})

	t.ObserveProperty(p, func(value interface{}, err error) {
		on := value.(bool)
		zlog.Trace().Str("name", info.Name).Bool("on", on).Msg("[main] Received update from thing device")
		acc.Switch.On.SetValue(on)
	})
	return &homeContext{accessory: acc, ConsumedThing: t}, acc.A
}

func (d homeContext) setCapability() {
}
