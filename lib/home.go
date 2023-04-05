package lib

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/service"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

type Context struct {
	*accessory.A
	Switch *service.Switch
}

func NewContext(info accessory.Info) *Context {
	a := Context{}
	a.A = accessory.New(info, accessory.TypeOther)
	a.Switch = service.NewSwitch()
	a.AddS(a.Switch.S)

	return &a
}

type homeContext struct {
	accessory *Context
	*consumer.ConsumedThing
}

func newHomeContext(info accessory.Info, t *consumer.ConsumedThing, p string) (*homeContext, *accessory.A) {
	acc := NewContext(info)
	data, err := t.ReadProperty(p)
	if err != nil {
		zlog.Error().Err(err).Msg("[main] Can't read context state")
	} else {
		state := data.(bool)
		zlog.Trace().Str("name", info.Name).Bool("value", state).Msg("[main] Set context initial value")
		acc.Switch.On.SetValue(state)
	}
	acc.Switch.On.OnValueRemoteUpdate(func(state bool) {
		zlog.Trace().Str("name", info.Name).Bool("on", state).Msg("[main] Received context update from Homekit")
		t.WriteProperty(p, state)
	})

	t.ObserveProperty(p, func(value interface{}, err error) {
		on := value.(bool)
		zlog.Trace().Str("name", info.Name).Bool("on", on).Msg("[main] Received context update from thing device")
		acc.Switch.On.SetValue(on)
	})
	return &homeContext{accessory: acc, ConsumedThing: t}, acc.A
}

func (d homeContext) setCapability() {
}
