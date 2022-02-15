package ipsetmgr

import (
	"sync"
)

// IPSetsManager provides an interface for managing multiple IP sets
type IPSetsManager struct {
	ipSets map[string]*IPSet
	lock   sync.RWMutex
}

// NewIPSetsManager will return a new IPSetsManager instance
func NewIPSetsManager() *IPSetsManager {
	return &IPSetsManager{
		ipSets: make(map[string]*IPSet),
	}
}

// TestIP tests a given ip against all sets in this set manager and returns arrays indicating
// which sets it matched against and which it did not
func (i *IPSetsManager) TestIP(ip string) ([]string, []string) {
	i.lock.RLock()

	positives := []string{}
	negatives := []string{}
	for ipSetName, ipSet := range i.ipSets {
		if ipSet.contains(ip) {
			positives = append(positives, ipSetName)
		} else {
			negatives = append(negatives, ipSetName)
		}
	}

	i.lock.RUnlock()
	return positives, negatives
}

// ReloadIPs will remove all sets from this manager, then add all the
// sets from the provided "loadedSets" map
func (i *IPSetsManager) ReloadIPs(loadedSets map[string][]string) {
	i.lock.Lock()

	// remove old sets
	i.ipSets = make(map[string]*IPSet)

	// populate new sets
	for name, ips := range loadedSets {
		newSet := NewIPSet()

		for _, ip := range ips {
			newSet.add(ip)
		}

		i.ipSets[name] = newSet
	}

	i.lock.Unlock()
}
