package mcfile

import (
	S "strings"
)

var excludeFilenamePrefixes = []string { ".", "_" }

var excludeFilepathContains = []string { ".." }

var excludeFilenameSuffixes = []string {
    "~", ".sh", ".rc", ".bashrc", "gtk", "gtr", "_echo", "_tkns", "_tree" }

// excludeFilenamepath returns true (plus a reason) for a file base
// name that matches blacklists for prefix or midfix or suffix. It
// does not check for absolute paths; this should be done elsehwere. 
func excludeFilenamepath(s string) (bool, string) {
     var reason string 
     for _, pfx := range excludeFilenamePrefixes {
     	 if S.HasPrefix(s, pfx) {
	    reason += "prefix<" + pfx + "> " 
	    }
	 }
     for _, sfx := range excludeFilenameSuffixes {
     	 if S.HasSuffix(s, sfx) {
	    reason += "suffix<" + sfx + "> " 
	    }
	 }
     for _, fpc := range excludeFilepathContains {
     	 if S.Contains(s, fpc) {
	    reason += "contains<" + fpc + "> " 
	    }
	 }
     return (reason != ""), reason 
}

