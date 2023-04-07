package lib

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

type shutterBasic struct {
	*consumer.ConsumedThing
	accessory *accessory.WindowCovering
}

func newShutterBasic(info accessory.Info, t *consumer.ConsumedThing) (*shutterBasic, *accessory.A) {
	acc := accessory.NewWindowCovering(info)
	initPosition(acc.WindowCovering, info, t, func(windowCovering *service.WindowCovering, t *consumer.ConsumedThing, value int) {
		if value == 100 {
			t.InvokeAction("open", nil)
			windowCovering.TargetPosition.SetValue(100)
		} else if value == 0 {
			t.InvokeAction("close", nil)
			windowCovering.TargetPosition.SetValue(0)
		}
	})

	t.ObserveProperty("open", func(value interface{}, err error) {
		var position int
		if value.(bool) {
			position = 100
		}
		zlog.Trace().Str("name", info.Name).Int("position", position).Msg("[main] Received shutter update from thing device")

		acc.WindowCovering.CurrentPosition.SetValue(position)
		acc.WindowCovering.TargetPosition.SetValue(position)
		acc.WindowCovering.PositionState.SetValue(characteristic.PositionStateStopped)
	})

	return &shutterBasic{accessory: acc, ConsumedThing: t}, acc.A
}

type shutterPosition struct {
	*consumer.ConsumedThing
	accessory *accessory.WindowCovering
}

func newShutterPosition(info accessory.Info, t *consumer.ConsumedThing) (*shutterPosition, *accessory.A) {
	acc := accessory.NewWindowCovering(info)
	initPosition(acc.WindowCovering, info, t, func(windowCovering *service.WindowCovering, t *consumer.ConsumedThing, value int) {
		if value == 100 {
			t.InvokeAction("open", nil)
			windowCovering.TargetPosition.SetValue(100)
		} else if value == 0 {
			t.InvokeAction("close", nil)
			windowCovering.TargetPosition.SetValue(0)
		} else {
			t.InvokeAction("setPosition", value)
			windowCovering.TargetPosition.SetValue(value)
		}
	})

	t.ObserveProperty("position", func(value interface{}, err error) {
		position := int(value.(float64))
		zlog.Trace().Str("name", info.Name).Int("position", position).Msg("[main] Received shutter update from thing device")

		acc.WindowCovering.CurrentPosition.SetValue(position)
		acc.WindowCovering.TargetPosition.SetValue(position) // set the TargetPosition to the CurrentPosition when the shutter is manually positioned. Otherwise the target != current position and the Home app will show that the shutter is moving.
		acc.WindowCovering.PositionState.SetValue(characteristic.PositionStateStopped)
	})

	return &shutterPosition{accessory: acc, ConsumedThing: t}, acc.A
}

func initPosition(windowCovering *service.WindowCovering, info accessory.Info, t *consumer.ConsumedThing, processValue func(*service.WindowCovering, *consumer.ConsumedThing, int)) {
	data, err := t.ReadProperty("position")
	if err != nil {
		zlog.Error().Err(err).Msg("[main] Can't read shutter position")
	} else {
		position := int(data.(float64))
		zlog.Trace().Str("name", info.Name).Int("value", position).Msg("[main] Set shutter initial value")

		windowCovering.CurrentPosition.SetValue(position)
		windowCovering.TargetPosition.SetValue(position)
		windowCovering.PositionState.SetValue(characteristic.PositionStateStopped)
	}

	windowCovering.TargetPosition.OnValueRemoteUpdate(func(value int) {
		zlog.Trace().Str("name", info.Name).Int("position", value).Msg("[main] Received shutter update from Homekit")
		current := windowCovering.CurrentPosition.Value()
		if value > current {
			windowCovering.PositionState.SetValue(characteristic.PositionStateIncreasing)
		} else if value < current {
			windowCovering.PositionState.SetValue(characteristic.PositionStateDecreasing)
		}
		processValue(windowCovering, t, value)
	})
}
