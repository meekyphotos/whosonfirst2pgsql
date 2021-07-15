package commands

import (
	"fmt"
	"regexp"

	"github.com/valyala/fastjson"
)

var latitudeRegex = regexp.MustCompile(`(?i).*:latitude`)
var longitudeRegex = regexp.MustCompile(`(?i).*:longitude`)
var preferredNames = regexp.MustCompile(`name:(?P<lang>.*)_x_preferred`)
var variantNames = regexp.MustCompile(`name:(?P<lang>.*)_x_variant`)
var countryCodeRegex = regexp.MustCompile(`.*:country`)

func fromJsonValue(value *fastjson.Value) Req {
	content := Req{}
	content.id = value.GetInt64("id")
	props := value.GetObject("properties")
	content.preferredNames = make(map[string]string)
	content.variantNames = make(map[string]string)
	content.metadata = make(map[string]string)
	content.continentId = value.GetInt64("properties", "wof:hierarchy", "continent_id")
	content.countryId = value.GetInt64("properties", "wof:hierarchy", "country_id")
	content.localityId = value.GetInt64("properties", "wof:hierarchy", "locality_id")
	content.macroregionId = value.GetInt64("properties", "wof:hierarchy", "macroregion_id")
	content.regionId = value.GetInt64("properties", "wof:hierarchy", "region_id")
	content.countyId = value.GetInt64("properties", "wof:hierarchy", "county_id")
	props.Visit(func(key []byte, v *fastjson.Value) {
		if latitudeRegex.Match(key) {
			if content.latitude == 0 {
				content.latitude = v.GetFloat64()
			}
		} else if longitudeRegex.Match(key) {
			if content.longitude == 0 {
				content.longitude = v.GetFloat64()
			}
		} else if countryCodeRegex.Match(key) {
			if content.countryCode == "" {
				content.countryCode = string(v.GetStringBytes())
			}
		} else {
			strKey := string(key)
			lang := preferredNames.FindStringSubmatch(strKey)
			if len(lang) > 0 {
				content.preferredNames[lang[1]] = string(v.GetArray()[0].GetStringBytes())
			} else {
				lang = variantNames.FindStringSubmatch(strKey)
				if len(lang) > 0 {
					content.variantNames[lang[1]] = string(v.GetArray()[0].GetStringBytes())
				} else if v.Type() == fastjson.TypeNumber {
					content.metadata[strKey] = fmt.Sprintf("%f", v.GetFloat64())
				} else if v.Type() == fastjson.TypeString {
					content.metadata[strKey] = string(v.GetStringBytes())
				}
			}
		}
	})
	return content
}
