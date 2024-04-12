package util

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"syscall"
	"unsafe"
)

type routeInfoInMemory struct {
	Family byte
	DstLen byte
	SrcLen byte
	TOS    byte

	Table    byte
	Protocol byte
	Scope    byte
	Type     byte

	Flags uint32
}

// rtInfo contains information on a single route.
type rtInfo struct {
	Src, Dst         *net.IPNet
	Gateway, PrefSrc net.IP
	// We currently ignore the InputIface.
	InputIface, OutputIface uint32
	Priority                uint32
}

// routeSlice implements sort.Interface to sort routes by Priority.
type routeSlice []*rtInfo

func (r routeSlice) Len() int {
	return len(r)
}
func (r routeSlice) Less(i, j int) bool {
	return r[i].Priority < r[j].Priority
}
func (r routeSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type ipAddrs struct {
	v4, v6 net.IP
}

type Router struct {
	ifaces map[int]net.Interface
	addrs  map[int]ipAddrs
	v4, v6 routeSlice
}

func NewRouter() (*Router, error) {
	rtr := &Router{}
	rtr.ifaces = make(map[int]net.Interface)
	rtr.addrs = make(map[int]ipAddrs)
	tab, err := syscall.NetlinkRIB(syscall.RTM_GETROUTE, syscall.AF_UNSPEC)
	if err != nil {
		return nil, err
	}
	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return nil, err
	}
loop:
	for _, m := range msgs {
		switch m.Header.Type {
		case syscall.NLMSG_DONE:
			break loop
		case syscall.RTM_NEWROUTE:
			rt := (*routeInfoInMemory)(unsafe.Pointer(&m.Data[0]))
			routeInfo := rtInfo{}
			attrs, err := syscall.ParseNetlinkRouteAttr(&m)
			if err != nil {
				return nil, err
			}
			if rt.Family != syscall.AF_INET && rt.Family != syscall.AF_INET6 {
				continue loop
			}
			for _, attr := range attrs {
				switch attr.Attr.Type {
				case syscall.RTA_DST:
					routeInfo.Dst = &net.IPNet{
						IP:   net.IP(attr.Value),
						Mask: net.CIDRMask(int(rt.DstLen), len(attr.Value)*8),
					}
				case syscall.RTA_SRC:
					routeInfo.Src = &net.IPNet{
						IP:   net.IP(attr.Value),
						Mask: net.CIDRMask(int(rt.SrcLen), len(attr.Value)*8),
					}
				case syscall.RTA_GATEWAY:
					routeInfo.Gateway = net.IP(attr.Value)
				case syscall.RTA_PREFSRC:
					routeInfo.PrefSrc = net.IP(attr.Value)
				case syscall.RTA_IIF:
					routeInfo.InputIface = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
				case syscall.RTA_OIF:
					routeInfo.OutputIface = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
				case syscall.RTA_PRIORITY:
					routeInfo.Priority = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
				}
			}
			if routeInfo.Dst == nil && routeInfo.Src == nil && routeInfo.Gateway == nil {
				continue loop
			}
			switch rt.Family {
			case syscall.AF_INET:
				rtr.v4 = append(rtr.v4, &routeInfo)
			case syscall.AF_INET6:
				rtr.v6 = append(rtr.v6, &routeInfo)
			default:
				// should not happen.
				continue loop
			}
		}
	}
	sort.Sort(rtr.v4)
	sort.Sort(rtr.v6)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		rtr.ifaces[iface.Index] = iface
		var addrs ipAddrs
		ifaceAddrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range ifaceAddrs {
			if inet, ok := addr.(*net.IPNet); ok {
				// Go has a nasty habit of giving you IPv4s as ::ffff:1.2.3.4 instead of 1.2.3.4.
				// We want to use mapped v4 addresses as v4 preferred addresses, never as v6
				// preferred addresses.
				if v4 := inet.IP.To4(); v4 != nil {
					if addrs.v4 == nil {
						addrs.v4 = v4
					}
				} else if addrs.v6 == nil {
					addrs.v6 = inet.IP
				}
			}
		}
		rtr.addrs[iface.Index] = addrs
	}
	return rtr, nil
}

func (r *Router) String() string {
	strs := []string{"ROUTER", "--- V4 ---"}
	for _, route := range r.v4 {
		strs = append(strs, fmt.Sprintf("%+v", *route))
	}
	strs = append(strs, "--- V6 ---")
	for _, route := range r.v6 {
		strs = append(strs, fmt.Sprintf("%+v", *route))
	}
	return strings.Join(strs, "\n")
}

func (r *Router) Route(dst net.IP) (iface net.Interface, gateway, preferredSrc net.IP, err error) {
	return r.RouteWithSrc(nil, nil, dst)
}

func (r *Router) RouteWithSrc(input net.HardwareAddr, src, dst net.IP) (iface net.Interface, gateway, preferredSrc net.IP, err error) {
	var ifaceIndex int
	switch {
	case dst.To4() != nil:
		ifaceIndex, gateway, preferredSrc, err = r.route(r.v4, input, src, dst)
	case dst.To16() != nil:
		ifaceIndex, gateway, preferredSrc, err = r.route(r.v6, input, src, dst)
	default:
		err = errors.New("IP is not valid as IPv4 or IPv6")
	}

	if err != nil {
		return
	}

	iface = r.ifaces[ifaceIndex]

	if preferredSrc == nil {
		switch {
		case dst.To4() != nil:
			preferredSrc = r.addrs[ifaceIndex].v4
		case dst.To16() != nil:
			preferredSrc = r.addrs[ifaceIndex].v6
		}
	}
	return
}

func (r *Router) route(routes routeSlice, input net.HardwareAddr, src, dst net.IP) (iface int, gateway, preferredSrc net.IP, err error) {
	var inputIndex uint32
	if input != nil {
		for i, iface := range r.ifaces {
			if bytes.Equal(input, iface.HardwareAddr) {
				inputIndex = uint32(i)
				break
			}
		}
	}
	var defaultGateway *rtInfo = nil
	for _, rt := range routes {
		if rt.InputIface != 0 && rt.InputIface != inputIndex {
			continue
		}
		if rt.Src == nil && rt.Dst == nil {
			defaultGateway = rt
			continue
		}
		if rt.Src != nil && !rt.Src.Contains(src) {
			continue
		}
		if rt.Dst != nil && !rt.Dst.Contains(dst) {
			continue
		}
		return int(rt.OutputIface), rt.Gateway, rt.PrefSrc, nil
	}

	if defaultGateway != nil {
		return int(defaultGateway.OutputIface), defaultGateway.Gateway, defaultGateway.PrefSrc, nil
	}
	err = fmt.Errorf("no route found for %v", dst)
	return
}
