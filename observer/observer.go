package observer

import (
	"net"
	"time"

	"github.com/mostlygeek/arp"
)

type AddressChange struct {
	Name string
	New  string
	Old  string
}

func Observe(config Config, state *State, events chan AddressChange) error {
	defer close(events)

	filters := map[string]interface{}{}
	for _, ip := range config.Filters {
		filters[ip] = nil
	}

	if config.Frequency == "" {
		*state = detectChanges(filters, *state, config.Aliases, events)
	}

	frequency, err := time.ParseDuration(config.Frequency)
	if err != nil {
		return err
	}

	for {
		*state = detectChanges(filters, *state, config.Aliases, events)
		<-time.After(frequency)
	}

	return nil
}

func getIPv4(addr net.Addr) string {
	if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
		if ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}

func detectChanges(filters map[string]interface{}, state State, aliases map[string]string, changes chan AddressChange) State {
	arpTable := arp.Table()

	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ip := getIPv4(addr); ip != "" {
				arpTable[ip] = iface.HardwareAddr.String()
			}
		}
	}

	for ip, newAddress := range arpTable {
		if _, found := filters[ip]; !found && len(filters) > 0 {
			continue
		}

		if oldAddress := state[ip]; oldAddress != newAddress {
			changes <- AddressChange{
				Name: aliasIP(ip, aliases),
				New:  newAddress,
				Old:  oldAddress,
			}
		}

		delete(state, ip)
	}

	for ip, oldAddress := range state {
		if _, found := filters[ip]; !found && len(filters) > 0 {
			continue
		}

		changes <- AddressChange{
			Name: aliasIP(ip, aliases),
			New:  "",
			Old:  oldAddress,
		}

	}

	return State(arpTable)
}

func aliasIP(ip string, aliases map[string]string) string {
	if alias := aliases[ip]; alias != "" {
		return alias
	}
	return ip
}
