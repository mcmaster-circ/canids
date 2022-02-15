package ipsetmgr

import (
	"fmt"
	"testing"
)

func compareTestIP(i *IPSetsManager, ip string, expPositives []string, expNegatives []string) error {
	positives, negatives := i.TestIP(ip)

	//TODO(Jon) implement this
	fmt.Printf("%s: %v %v =?= %v %v\n", ip, positives, negatives, expPositives, expNegatives)

	return nil
}

func TestAddAndContains(t *testing.T) {
	ipSetsMgr := NewIPSetsManager()
	ipSetsMgr.ReloadIPs(map[string][]string{
		"testset1": []string{"10.0.0.53"},
	})

	if err := compareTestIP(ipSetsMgr, "10.0.0.53", []string{"testset1"}, []string{}); err != nil {
		t.Error(err)
	}
	if err := compareTestIP(ipSetsMgr, "10.0.0.12", []string{}, []string{"testset1"}); err != nil {
		t.Error(err)
	}

	ipSetsMgr.ReloadIPs(map[string][]string{
		"testset1": []string{"10.0.0.53", "10.0.0.12"},
	})

	if err := compareTestIP(ipSetsMgr, "10.0.0.12", []string{"testset1"}, []string{}); err != nil {
		t.Error(err)
	}
}

func TestAddAndContainsForRanges(t *testing.T) {
	ipSetsMgr := NewIPSetsManager()
	ipSetsMgr.ReloadIPs(map[string][]string{
		"testset1": []string{"10.0.1.0/24"},
	})

	if err := compareTestIP(ipSetsMgr, "10.0.1.53", []string{"testset1"}, []string{}); err != nil {
		t.Error(err)
	}
	if err := compareTestIP(ipSetsMgr, "10.0.0.53", []string{}, []string{"testset1"}); err != nil {
		t.Error(err)
	}
}
