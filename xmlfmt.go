////////////////////////////////////////////////////////////////////////////
// Porgram: xml-beautify-regexp.go
// Purpose: Go XML Beautify from XML string using pure regexp
// Authors: Antonio Sun (c) 2016, All rights reserved
// Credits: diotalevi http://www.perlmonks.org/?node_id=261292
////////////////////////////////////////////////////////////////////////////

package xmlfmt

import (
	"regexp"
	"strings"
)

var (
	reg = regexp.MustCompile(`<(/?)([^>]+)(/?)>`)
	// NL is the newline string used in XML ouput, define for DOS-convenient.
	NL = "\r\n"
)

// FormatXML will format the XML string in a readable way
func FormatXML(xmls, prefix, indent string) string {
	src := regexp.MustCompile(`>\s+<`).ReplaceAllString(xmls, "><")

	rf := replaceTag(prefix, indent)
	return (prefix + reg.ReplaceAllStringFunc(src, rf))
}

// replaceTag returns a closure function to do 's/(?<=>)\s+(?=<)//g; s(<(/?)([^>]+)(/?)>)($indent+=$3?0:$1?-1:1;"<$1$2$3>"."\n".("  "x$indent))ge' as in Perl
func replaceTag(prefix, indent string) func(string) string {
	indentLevel := 0
	return func(m string) string {
		parts := reg.FindStringSubmatch(m)
		// $3: A <foo/> tag. No alteration to indentation.
		// $1: A closing </foo> tag. Drop one indentation level
		// else: An opening <foo> tag. Increase one indentation level
		if len(parts[3]) == 0 && len(parts[1]) != 0 {
			indentLevel--
		} else {
			indentLevel++
		}
		return "<" + parts[1] + parts[2] + parts[3] + ">" +
			NL + prefix + strings.Repeat(indent, indentLevel)
	}
}