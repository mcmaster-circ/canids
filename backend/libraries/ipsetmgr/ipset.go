package ipsetmgr

import (
	"net"
	"strings"

	"github.com/yl2chen/cidranger"
)

// IPSet lets you check if an exact ip is contained within a set of ips.
type IPSet struct {
	ranger    *cidranger.Ranger
	singleIPs map[string]bool
}

// NewIPSet returns an instance of a IPSet
func NewIPSet() *IPSet {
	ranger := cidranger.NewPCTrieRanger()
	newIPSet := &IPSet{
		ranger:    &ranger,
		singleIPs: make(map[string]bool),
	}
	return newIPSet
}

// contains will return true if the given ip exists within a set
func (s *IPSet) contains(ip string) bool {
	if s.singleIPs[ip] {
		return true
	}
	r := *s.ranger
	contains, err := r.Contains(net.ParseIP(ip))
	if err != nil {
		return false
	}
	return contains
}

// add will add a new ip or ip range to the set
func (s *IPSet) add(ip string) {
	if strings.Contains(ip, "/") {
		s.addRangedIP(ip)
	} else {
		s.addSingleIP(ip)
	}
}

// addSingleIP will add a singular ip to the set
func (s *IPSet) addSingleIP(ip string) {
	s.singleIPs[ip] = true
}

// add rangedIP will add an ip range to the set
func (s *IPSet) addRangedIP(ip string) {
	_, network, err := net.ParseCIDR(ip)
	if err != nil {
		panic(err)
	}
	r := *s.ranger
	r.Insert(cidranger.NewBasicRangerEntry(*network))
}

// removeAll will remove all the ips in the set
func (s *IPSet) removeAll() {
	ranger := cidranger.NewPCTrieRanger()
	s.ranger = &ranger
}
