// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package websocket provides websocket streaming ingestion services.
package websocket

import (
	"fmt"
	"net"

	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// geoIPASN returns the ASN of a given IP address. If the ASN cannot be
// determined, an empty string will be returned.
func geoIPASN(s *state.State, ipAddress string) string {
	ip := net.ParseIP(ipAddress)
	result, err := s.GeoIPASN.ASN(ip)
	if err != nil {
		return ""
	}
	return "AS" + fmt.Sprint(result.AutonomousSystemNumber) + " " + result.AutonomousSystemOrganization
}

// geoIPCity returns the city of a given IP address. If the coty cannot be
// determined, an empty string will be returned.
func geoIPCity(s *state.State, ipAddress string) string {
	ip := net.ParseIP(ipAddress)
	result, err := s.GeoIPCity.City(ip)
	if err != nil {
		return ""
	}
	city, ok := result.City.Names["en"]
	if !ok {
		return ""
	}
	return city
}

// geoIPCountry returns the country of a given IP address. If the country cannot
// be determined, an empty string will be returned.
func geoIPCountry(s *state.State, ipAddress string) string {
	ip := net.ParseIP(ipAddress)
	result, err := s.GeoIPCountry.Country(ip)
	if err != nil {
		return ""
	}
	country, ok := result.Country.Names["en"]
	if !ok {
		return ""
	}
	return country
}
