package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gateway-homekit/lib"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	eria "github.com/project-eria/eria-core"
	zlog "github.com/rs/zerolog/log"
)

var config = struct {
	Pin         string        `yaml:"pin" required:"true"`
	StoragePath string        `yaml:"storagePath" required:"true"`
	Devices     []ThingDevice `yaml:"devices" required:"true"`
}{}

// ThingDevice represents a Thing Device connection
type ThingDevice struct {
	Url      string `yaml:"url" required:"true"`
	Type     string `yaml:"type" required:"true"`
	Property string `yaml:"property"`
}

func init() {
	eria.Init("ERIA Homekit Gateway")
}

func main() {
	defer func() {
		zlog.Info().Msg("[main] Stopped")
	}()

	//	dnslog.Debug.Enable()

	// Loading config
	eria.LoadConfig(&config)

	// Store the data in the "./db" directory.
	fs := hap.NewFsStore("./db")

	bridge := accessory.NewBridge(accessory.Info{
		Name:         "homekit-gateway",
		Manufacturer: "ERIA",
	})

	accessories := setup()

	// Create the hap server.
	server, err := hap.NewServer(fs, bridge.A, accessories...)
	if err != nil {
		zlog.Fatal().Err(err).Msg("[main] New HAP Server")
	}
	server.Pin = config.Pin

	// Setup a listener for interrupts and SIGTERM signals
	// to stop the server.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		// Stop delivering signals.
		signal.Stop(c)
		// Cancel the context to stop the server.
		cancel()
	}()

	// Run the server.
	err = server.ListenAndServe(ctx)
	if err != nil {
		zlog.Fatal().Err(err).Msg("[main] ListenAndServe")
	}
	// eria.WaitForSignal()
}

// setup : create devices, register ...
func setup() []*accessory.A {
	accessories := []*accessory.A{}
	eriaClient := eria.NewClient()

	for _, device := range config.Devices {
		//		zlog.Info().Str("label", device.Label).Str("url", device.Url).Msg("[main] Connecting remote Things")
		remoteThing, err := eriaClient.ConnectThing(device.Url)
		if err == nil {
			acc, err := lib.NewAccessory(remoteThing, device.Type, device.Property)
			if err != nil {
				zlog.Error().Err(err).Msg("[main] initAccessory")
			}
			accessories = append(accessories, acc)
		} else {
			zlog.Error().Err(err).Msg("[main] Can't connect remote Thing")
		}
	}
	return accessories
}
