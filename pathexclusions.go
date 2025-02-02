package mcfile

import (
	S "strings"
)

var excludePrefixes = []string { ".", "_" }

var excludeContains = []string { ".." }

var excludeSuffixes = []string {
    "~", ".sh", ".rc", ".bashrc", "gtk", "gtr", "_echo", "_tkns", "_tree" }

// excludeFilenamepath returns true (plus a reason) for a file base
// name that matches a blacklist for prefix or "midfix" or suffix.
// Excluded prefixes are also checked for following a path separator.
func excludeFilenamepath(s string) (bool, string) {
     var reason string  
     for _, pfx := range excludePrefixes {
     	 if S.HasPrefix(s, pfx) {
	    reason += "prefix<" + pfx + "> " 
	    }
	 }
     for _, sfx := range excludeSuffixes {
     	 if S.HasSuffix(s, sfx) {
	    reason += "suffix<" + sfx + "> " 
	    }
	 }
     for _, fpc := range excludeContains {
     	 if S.Contains(s, fpc) {
	    reason += "contains<" + fpc + "> " 
	    }
	 }
     for _, pfx := range excludePrefixes {
     	 if S.Contains(s, "/" + pfx) {
	    reason += "/+prefix<" + pfx + "> " 
	    }
	 }
     return (reason != ""), reason 
}

