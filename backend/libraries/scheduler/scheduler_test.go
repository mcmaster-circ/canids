package scheduler

import (
	"fmt"
	"testing"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ipsetmgr"
)

func TestProvision(t *testing.T) {
	ipSetsMgr := ipsetmgr.NewIPSetsManager()

	ProvisionOnce(map[string]string{
		"firehol_abusers_1d":  "https://iplists.firehol.org/files/firehol_abusers_1d.netset",
		"firehol_abusers_30d": "https://iplists.firehol.org/files/firehol_abusers_30d.netset",
		"firehol_anonymous":   "https://iplists.firehol.org/files/firehol_anonymous.netset",
		"firehol_level1":      "https://iplists.firehol.org/files/firehol_level1.netset",
		"firehol_level2":      "https://iplists.firehol.org/files/firehol_level2.netset",
		"firehol_level3":      "https://iplists.firehol.org/files/firehol_level3.netset",
	}, ipSetsMgr)

	ipsToTest := []string{
		"10.189.34.102", "10.189.34.112", "10.189.33.7",
		"10.189.34.150", "fe80::5df2:8e94:cd04:319a", "10.189.34.15",
		"192.168.2.81", "10.189.33.47", "10.189.32.15", "10.189.34.124",
	}

	for _, ip := range ipsToTest {
		positives, negatives := ipSetsMgr.TestIP(ip)
		fmt.Printf("%s | positives %v | negatives %v\n", ip, positives, negatives)
	}

	fmt.Println("testing: 1.1.222.168")
	positives, negatives := ipSetsMgr.TestIP("1.1.222.168")
	fmt.Println("positives:")
	fmt.Println(positives)
	fmt.Println("negatives:")
	fmt.Println(negatives)

	fmt.Println("testing: 1.10.16.0")
	positives, negatives = ipSetsMgr.TestIP("1.10.16.0")
	fmt.Println("positives:")
	fmt.Println(positives)
	fmt.Println("negatives:")
	fmt.Println(negatives)
}
