package scheduler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ipsetmgr"
)

// Provision will accept: a map of ip set names to their urls, a time interval to schedule provisioning,
// and an IPSetsManager instance. It will regularly provision the ipsetmgr with the contents of the url
// based on the given time interval.
func Provision(
	urls map[string]string,
	waitTime time.Duration,
	ipSetsMgr *ipsetmgr.IPSetsManager,
) {
	ticker := time.NewTicker(waitTime)

	// ping google to see if we are on the internet, if not dont load the ip sets
	_, err := http.Get("https://www.google.com/")
	if err == nil {
		fmt.Printf("Ping to google succeeded, we're on the internet. Loading alarm IP sets.\n")
	} else {
		fmt.Printf("Ping to google failed. Do not load alarm IP sets.\n")
		return
	}

	// do initial provision, this takes a while (~3 mins sometimes)
	err = ProvisionOnce(urls, ipSetsMgr)
	if err != nil {
		fmt.Printf("Error provisioning alarm: %d", err)
	}

	// start loop that does periodic refreshes
	go func() {
		for {
			select {
			case <-ticker.C:
				err := ProvisionOnce(urls, ipSetsMgr)
				if err != nil {
					fmt.Printf("Error provisioning alarm: %d", err)
				}
			}
		}
	}()
}

// ProvisionOnce will iterate through the given urls and store the retrieved ips into the ip set manager.
func ProvisionOnce(
	urls map[string]string,
	ipSetsMgr *ipsetmgr.IPSetsManager,
) error {
	t0 := time.Now()
	loadedSets := make(map[string][]string)
	for name, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		loadedSets[name] = getIPsFromText(string(bodyBytes))
	}
	fmt.Printf("Loaded set queries: %d ms\n", time.Now().Sub(t0).Milliseconds())

	t0 = time.Now()
	ipSetsMgr.ReloadIPs(loadedSets)
	fmt.Printf("Update ip set manager: %d ms\n", time.Now().Sub(t0).Milliseconds())

	return nil
}

// getIPsFromText will parse through a string text and return the list of IPs,
// ignoring lines that start with "#".
func getIPsFromText(text string) []string {

	textLines := strings.Split(text, "\n")
	results := make([]string, 0)
	for _, line := range textLines {
		if len(line) > 0 && string(line[0]) == "#" {
			continue
		}
		results = append(results, line)
	}
	return results
}
