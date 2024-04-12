package wg

import (
	"net"

	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
)

func (w *WgIface) createWithUserspace() error {
	tunIface, err := tun.CreateTUN(w.Name, w.MTU)
	if err != nil {
		return err
	}

	w.Interface = tunIface

	// We need to create a wireguard-go device and listen to configuration requests
	w.tunDevice = device.NewDevice(tunIface, conn.NewDefaultBind(), device.NewLogger(device.LogLevelSilent, "[fahi] "))

	err = w.assignAddr()
	if err != nil {
		return err
	}

	w.uapiListener, err = getUAPI(w.Name)
	if err != nil {
		return err
	}

	go func(uapi net.Listener) {
		for {
			uapiConn, uapiErr := uapi.Accept()
			if uapiErr != nil {
				log.Traceln("uapi Accept failed with error: ", uapiErr)
				return
			}
			go w.tunDevice.IpcHandle(uapiConn)
		}
	}(w.uapiListener)

	log.Debugln("UAPI listener started")

	err = w.configure()
	if err != nil {
		return err
	}

	err = w.tunDevice.Up()
	if err != nil {
		return err
	}

	return nil
}

// getUAPI returns a Listener
func getUAPI(iface string) (net.Listener, error) {
	tunSock, err := ipc.UAPIOpen(iface)
	if err != nil {
		return nil, err
	}
	return ipc.UAPIListen(iface, tunSock)
}
