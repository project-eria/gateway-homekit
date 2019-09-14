package main

import (
	"os"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/project-eria/eria-base"
	logger "github.com/project-eria/eria-logger"
	"github.com/project-eria/xaal-go"
	"github.com/project-eria/xaal-go/device"
	"github.com/project-eria/xaal-go/schemas"
)

var (
	// Version is a placeholder that will receive the git tag version during build time
	Version = "-"
)

const configFile = "gateway-homekit.json"

func setupDev(dev *device.Device) {
	dev.VendorID = "ERIA"
	dev.ProductID = "Apple Homekit Gateway"
	dev.Info = "gateway.homekit"
	dev.Version = Version
}

var config = struct {
	GWXaalAddr  string
	Pin         string `required:"true"`
	StoragePath string `required:"true"`
	Devices     []configDevice
}{}

type configDevice struct {
	Type     string
	Name     string
	XaalAddr string
}

var (
	_gw            *device.Device
	_devicesByXAAL map[string]*configDevice
)

func main() {
	defer os.Exit(0)

	eria.AddShowVersion(Version)

	logger.Module("main").Infof("Starting gateway-homekit %s...", Version)

	// Loading config
	cm := eria.LoadConfig(configFile, &config)
	defer cm.Close()

	// Init xAAL engine
	eria.InitEngine()

	// TODO RUN IN TRACE MODE	dnslog.Debug.Enable()

	accessories := setup()

	// Save for new Address during setup
	cm.Save()

	xaal.AddRxHandler(updateFromXAAL)

	t, err := hc.NewIPTransport(hc.Config{Pin: config.Pin, StoragePath: config.StoragePath}, newBridgeAccessory(), accessories...)
	if err != nil {
		logger.Module("main").Fatal(err)
	}

	go xaal.Run()
	defer xaal.Stop()

	hc.OnTermination(func() {
		<-t.Stop()
	})

	go t.Start()

	// Configure the schedulers
	eria.WaitForExit()
}

// setup : create devices, register ...
func setup() []*accessory.Accessory {
	// gw
	_gw, _ = schemas.Gateway(config.GWXaalAddr)
	setupDev(_gw)
	xaal.AddDevice(_gw)

	_devicesByXAAL = map[string]*configDevice{}
	accessories := []*accessory.Accessory{}

	for i := range config.Devices {
		confDev := &config.Devices[i]
		_devicesByXAAL[confDev.XaalAddr] = confDev
		accessory, err := initAccessory(confDev)
		if err != nil {
			logger.Module("main").Error(err)
		}
		accessories = append(accessories, accessory)
	}
	return accessories
}

func newBridgeAccessory() *accessory.Accessory {
	info := accessory.Info{
		Name:         "xaal-homekit-gateway",
		Manufacturer: "ERIA",
	}
	return accessory.New(info, accessory.TypeBridge)
}
