////////////////////////////////////////////////////////////////////////////
// Porgram: xmlfmt.go
// Purpose: Go XML Beautify from XML string using pure string manipulation
// Authors: Antonio Sun (c) 2016-2019, All rights reserved
////////////////////////////////////////////////////////////////////////////

package xmlfmt

import (
	"regexp"
	"strings"
)

var (
	reg = regexp.MustCompile(`<([/!]?)([^>]+?)(/?)>`)
	// NL is the newline string used in XML output
	NL = "\n"
)

// FormatXML will (purly) reformat the XML string in a readable way, without any rewriting/altering the structure
func FormatXML(xmls, prefix, indent string) string {
	// replace all whitespace between tags
	src := regexp.MustCompile(`(?s)>\s+<`).ReplaceAllString(xmls, "><")

	rf := replaceTag(prefix, indent)
	return (prefix + reg.ReplaceAllStringFunc(src, rf))
}

// replaceTag returns a closure function to do 's/(?<=>)\s+(?=<)//g; s(<(/?)([^>]+?)(/?)>)($indent+=$3?0:$1?-1:1;"<$1$2$3>"."\n".("  "x$indent))ge' as in Perl
// and deal with comments as well
func replaceTag(prefix, indent string) func(string) string {
	indentLevel := 0
	justOpened := false
	return func(m string) string {
		// head elem
		if strings.HasPrefix(m, "<?xml") {
			justOpened = false
			return prefix + strings.Repeat(indent, indentLevel) + m + NL
		}

		// empty elem
		if strings.HasSuffix(m, "/>") {
			if justOpened {
				justOpened = false
				return NL + prefix + strings.Repeat(indent, indentLevel) + m + NL
			} else {
				return prefix + strings.Repeat(indent, indentLevel) + m + NL
			}
		}
		// comment elem
		if strings.HasPrefix(m, "<!") {
			justOpened = false
			return NL + prefix + strings.Repeat(indent, indentLevel) + m
		}
		// end elem
		if strings.HasPrefix(m, "</") {
			defer func() {
				justOpened = false
			}()
			indentLevel--
			if justOpened {
				return m + NL
			}
			return prefix + strings.Repeat(indent, indentLevel) + m + NL
		}

		defer func() {
			indentLevel++
		}()

		if justOpened {
			// indentLevel++
			return NL + prefix + strings.Repeat(indent, indentLevel) + m
		}

		justOpened = true
		return prefix + strings.Repeat(indent, indentLevel) + m
	}
}
