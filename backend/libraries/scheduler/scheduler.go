package scheduler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ipsetmgr"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/uuid"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// Provision will accept: a map of ip set names to their urls, a time interval to schedule provisioning,
// and an IPSetsManager instance. It will regularly provision the ipsetmgr with the contents of the url
// based on the given time interval.
func Provision(
	s *state.State,
	waitTime time.Duration,
	ipSetsMgr *ipsetmgr.IPSetsManager,
) error {
	s.Log.Info("[scheduler] provisioning ip sets")
	ticker := time.NewTicker(waitTime)

	urls := make(map[string]string)

	// check if blacklist index exists, if not create it
	exists, err := s.Elastic.IndexExists("blacklist").Do(s.ElasticCtx)
	if err != nil {
		return err
	}
	if exists {
		// load the urls from the index
		s.Log.Info("[scheduler] loading blacklist index")
		urls = loadBlacklists(s)
	} else {
		// create a new index for blacklists
		s.Log.Info("[scheduler] creating blacklist index")
		urls = createAndLoadDefaultBlacklists(s)
	}

	// ping google to see if we are on the internet, if not dont load the ip sets
	_, err = http.Get("https://www.google.com/")
	if err == nil {
		fmt.Printf("Ping to google succeeded, we're on the internet. Loading alarm IP sets.\n")
	} else {
		fmt.Printf("Ping to google failed. Do not load alarm IP sets.\n")
		return nil
	}

	// do initial provision, this takes a while (~3 mins sometimes)
	fmt.Println("Provisioning alarm IP sets...")
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

	return nil
}

// ProvisionOnce will iterate through the given urls and store the retrieved ips into the ip set manager.
func ProvisionOnce(
	urls map[string]string,
	ipSetsMgr *ipsetmgr.IPSetsManager,
) error {
	t0 := time.Now()
	print("\n\n\n\n")
	loadedSets := make(map[string][]string)
	for name, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		println(url)
		startBadTime := time.Now()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("Elapsed time for above url: %d ms \n\n", time.Now().Sub(startBadTime).Milliseconds())
		if err != nil {
			return err
		}

		loadedSets[name] = getIPsFromText(string(bodyBytes))
	}
	fmt.Printf("Loaded set queries: %d ms\n\n\n\n", time.Now().Sub(t0).Milliseconds())

	t0 = time.Now()
	ipSetsMgr.ReloadIPs(loadedSets)
	fmt.Printf("Update ip set manager: %d ms\n", time.Now().Sub(t0).Milliseconds())

	return nil
}

// Refresh will refresh the ip sets in the alarm manager
func Refresh(s *state.State) {
	blacklistMap := loadBlacklists(s)

	fmt.Println("Refreshing alarm IP sets...")

	err := ProvisionOnce(blacklistMap, s.AlarmManager)
	if err != nil {
		fmt.Printf("Error refreshing ip list: %d", err)
	}
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

// LoadBlacklists will load the blacklists from the database into a map
func loadBlacklists(s *state.State) map[string]string {
	blacklists, err := elasticsearch.AllBlacklists(s)
	if err != nil {
		s.Log.Error("error getting all blacklists ", err)
		return nil
	}

	// convert blacklist documents to map
	blacklistMap := make(map[string]string)
	for _, blacklist := range blacklists {
		blacklistMap[blacklist.Name] = blacklist.URL
	}

	return blacklistMap
}

// CreateAndLoadDefaultBlacklists will create an index with the default blacklists
func createAndLoadDefaultBlacklists(s *state.State) map[string]string {
	blacklistMap := map[string]string{
		"firehol_abusers_1d":  "https://iplists.firehol.org/files/firehol_abusers_1d.netset",
		"firehol_abusers_30d": "https://iplists.firehol.org/files/firehol_abusers_30d.netset",
		"firehol_anonymous":   "https://iplists.firehol.org/files/firehol_anonymous.netset",
		"firehol_level1":      "https://iplists.firehol.org/files/firehol_level1.netset",
		"firehol_level2":      "https://iplists.firehol.org/files/firehol_level2.netset",
		"firehol_level3":      "https://iplists.firehol.org/files/firehol_level3.netset",
	}

	s.Elastic.CreateIndex("blacklist").Do(s.ElasticCtx)

	for name, url := range blacklistMap {
		blacklist := elasticsearch.DocumentBlacklist{
			UUID: uuid.Generate(),
			Name: name,
			URL:  url,
		}
		_, err := blacklist.Index(s)
		if err != nil {
			s.Log.Error("error indexing new blacklist ", err)
			return nil
		}
	}

	return blacklistMap
}
