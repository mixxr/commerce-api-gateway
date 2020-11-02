package utils

import (
	"regexp"
	"strings"
)

type Sanitizer struct {
	regExpUrlPath   *regexp.Regexp
	regExpUrlSearch *regexp.Regexp
	regExpUrlNumber *regexp.Regexp
	regExpDesc      *regexp.Regexp
	regExpTags      *regexp.Regexp
	regExpTagsDelim *regexp.Regexp
}

func NewSanitizer() *Sanitizer {
	o := Sanitizer{
		regExpUrlPath:   regexp.MustCompile("^[a-z0-9_]{1,32}$"),
		regExpUrlSearch: regexp.MustCompile("^[a-z0-9_]{0,32}-?$"), // search token like "servi-"
		regExpUrlNumber: regexp.MustCompile("^0-9]{1,19}$"),        // int64
		regExpDesc:      regexp.MustCompile("[^.,\\- a-zA-Z0-9]+"),
		regExpTags:      regexp.MustCompile("[^\\-a-zA-Z0-9]+"),
		regExpTagsDelim: regexp.MustCompile("[/]{2,}"), // duplicate /
	}
	return &o
}

// GetToken returns a max 32-length string in lowercase that comply to [a-z0-9_]{1,32}
func (o *Sanitizer) GetToken(urlPath string) string {
	if len(urlPath) > 0 {
		urlPath = strings.ToLower(urlPath)
		urlPath = o.regExpUrlPath.FindString(urlPath)
	}

	return urlPath
}

// GetSearchToken returns a max 32-length string in lowercase that comply to [a-z0-9_]{0,32}-?
func (o *Sanitizer) GetSearchToken(urlPath string) string {
	if len(urlPath) > 0 {
		urlPath = strings.ToLower(urlPath)
		urlPath = o.regExpUrlSearch.FindString(urlPath)
	}

	return urlPath
}

func Trunc(str string, maxLen int, tail string) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen-len(tail)] + tail
}

// GetDesc returns a max 256-length string with no special characters like $%#@!()*ˆ?/!';"[]
func (o *Sanitizer) GetDescr(desc string) string {
	if len(desc) > 0 {
		desc = o.regExpDesc.ReplaceAllString(desc, "")
		return Trunc(desc, 256, "...")
	}

	return desc
}

// GetTags converts tag1/tag2/... format in tag1,tag2,... format, max 256-length, no special characters like $%#@!()*ˆ?/!';"[]
func (o *Sanitizer) GetTags(tagsstr string) string {
	tagsstr = strings.TrimRight(tagsstr, "/")
	tagsstr = strings.TrimLeft(tagsstr, "/")
	if len(tagsstr) <= 1 {
		return ""
	}
	tagsstr = o.regExpTagsDelim.ReplaceAllString(tagsstr, "/") // replacing duplicate /
	tags := strings.Split(tagsstr, "/")

	for i, tag := range tags {
		tags[i] = o.regExpTags.ReplaceAllString(tag, "")
	}

	return Trunc(strings.Join(tags, ","), 256, "")
}

// CheckToken returns the string in lowercase if it does not contain special characters otherwise an empty string
func (o *Sanitizer) CheckToken(urlPath string) string {
	if o.regExpUrlPath.MatchString(urlPath) {
		return strings.ToLower(urlPath)
	}

	return ""
}

func (o *Sanitizer) CheckTokens(owner string, service string) (string, string) {
	if o.regExpUrlPath.MatchString(owner) && o.regExpUrlPath.MatchString(service) {
		return strings.ToLower(owner), strings.ToLower(service)
	}

	return "", ""
}

// CheckNumber returns the string if it is a number otherwise an empty string
func (o *Sanitizer) CheckNumber(urlPath string) string {
	if o.regExpUrlNumber.MatchString(urlPath) {
		return urlPath
	}
	return ""
}
