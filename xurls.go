/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

import "regexp"

//go:generate go run tools/tldsgen/main.go
//go:generate go run tools/regexgen/main.go

const (
	letters   = "a-zA-Z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF"
	iriChar   = letters + `0-9`
	pathChar  = iriChar + `/\-+_@&=#$~*%.,:;'"()?!`
	endChar   = iriChar + `/\-+_@&=#$~*%`
	ipv4Addr  = `(25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[1-9])\.(25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[0-9])`
	ipv6Addr  = `([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:[0-9a-fA-F]{0,4}|:[0-9a-fA-F]{1,4})?|(:[0-9a-fA-F]{1,4}){0,2})|(:[0-9a-fA-F]{1,4}){0,3})|(:[0-9a-fA-F]{1,4}){0,4})|:(:[0-9a-fA-F]{1,4}){0,5})((:[0-9a-fA-F]{1,4}){2}|:(25[0-5]|(2[0-4]|1[0-9]|[1-9])?[0-9])(\.(25[0-5]|(2[0-4]|1[0-9]|[1-9])?[0-9])){3})|(([0-9a-fA-F]{1,4}:){1,6}|:):[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){7}:`
	ipAddr    = `(` + ipv4Addr + `|` + ipv6Addr + `)`
	iri       = `[` + iriChar + `]([` + iriChar + `\-]*[` + iriChar + `])?`
	domain    = `(` + iri + `\.)+`
	hostName  = `(` + domain + gtld + `|` + ipAddr + `)`
	wellParen = `([` + pathChar + `]*(\([` + pathChar + `]*\))+)+`
	pathCont  = `(` + wellParen + `|[` + pathChar + `]*[` + endChar + `])`
	path      =  `(/` + pathCont + `?|\b|$)`
	webURL    = hostName + `(:[0-9]{1,5})?` + path
	email     = `[a-zA-Z0-9._%\-+]+@` + hostName

	commonScheme = `[a-zA-Z.\-+]+://`
	scheme       = `(` + commonScheme + `|` + otherScheme + `)`
	strict       = `(\b|^)` + scheme + pathCont
	relaxed      = strict + `|` + webURL + `|` + email
)

var (
	// Relaxed matches all the urls it can find
	Relaxed = regexp.MustCompile(relaxed)
	// Strict only matches urls with a scheme to avoid false positives
	Strict = regexp.MustCompile(strict)
)

func init() {
	Relaxed.Longest()
	Strict.Longest()
}

func StrictMatching(schemeExp string) *regexp.Regexp {
	strictMatching := `(\b|^)(` + schemeExp + `)` + pathCont
	re := regexp.MustCompile(strictMatching)
	re.Longest()
	return re
}
