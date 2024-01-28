package iplocation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type GeoIpLocation struct {
	Ip               string  `json:"ip"`
	Latitude         string  `json:"latitude" bson:"latitude"`
	Longitude        string  `json:"longitude" bson:"longitude"`
	City             *string `json:"city" bson:"city"`
	Region           *string `json:"region" bson:"state"`
	Country          string  `json:"country" bson:"country"`
	OrganizationName string  `json:"organization_name" bson:"organization_name"`
}

func GetGeoLocationByIp(ip string) (loc *GeoIpLocation, err error) {
	if os.Getenv("environment") == "dev" {
		ip = "66.241.125.71"
	}
	url := fmt.Sprintf("https://get.geojs.io/v1/ip/geo/%s.json", ip)
	resp, err := http.Get(url)
	if err != nil {
		zap.L().Error("error get geo location", zap.Error(err), zap.String("url", url))
		return nil, err
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&loc)
	zap.L().Info("get geo location", zap.String("ip", ip), zap.Any("location", loc))
	return
}
