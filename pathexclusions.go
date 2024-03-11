package mcfile

import (
	S "strings"
)

var excludeFilenamePrefixes = []string { ".", "_" }

var excludeFilepathContains = []string { ".." }

var excludeFilenameSuffixes = []string { "~", ".sh", ".rc", ".bashrc" }

// excludeFilenamepath does not check for absolute paths;
// this should be done elsehwere. 
func excludeFilenamepath(s string) (bool, string) {
     var reason string 
     for _, pfx := range excludeFilenamePrefixes {
     	 if S.HasPrefix(s, pfx) {
	    reason += "prefix "
	    }
	 }
     for _, sfx := range excludeFilenameSuffixes {
     	 if S.HasSuffix(s, sfx) {
	    reason += "suffix "
	    }
	 }
     for _, fpc := range excludeFilepathContains {
     	 if S.Contains(s, fpc) {
	    reason += "contains "
	    }
	 }
     return (reason != ""), reason 
}

