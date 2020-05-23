package mcfile

import (
	// "flag"
	flag "github.com/spf13/pflag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	FP "path/filepath"

	"github.com/fbaube/db"
	FU "github.com/fbaube/fileutils"
	"github.com/fbaube/gparse"
	SU "github.com/fbaube/stringutils"
	MU "github.com/fbaube/miscutils"
	WU "github.com/fbaube/wasmutils"
	"errors"
)

// Using this messes up the default Usage(), so **avoid** `flag.FlagSet`
// var fs flag.FlagSet

// ConfigurationArguments can probably be used with various 3rd-party utilities.
type ConfigurationArguments struct {
	AppName   string
	DBdirPath string
	DBhandle  *db.MmmcDB
	In, Out, XmlCat,
	XmlCatSearch FU.BasicPath // NOT ptr! Barfs at startup.
	RestPort int
	// CLI flags
	FollowSymLinks, Pritt, DBdoImport, Help, Debug, GroupGenerated, Validate, DBdoZeroOut bool
	// Result of processing CLI arg for input file(s)
	SingleFile bool
	// Result of processing CLI args (-c, -s)
	*gparse.XmlCatalog
}

// var InFP, OutFP, XmlCatFP, XmlCatSearchFP string
var xmlCatalogs []*gparse.XmlCatalog

// CA maybe should not be exported.
var CA ConfigurationArguments
// var rootNode, currentNode *gtree.GTag

func myUsage() {
	//  Println(CA.AppName, "[-d] [-g] [-h] [-m] [-p] [-v] [-z] [-D] [-o outfile] [-d dbdir] Infile")
	fmt.Println(CA.AppName, "[-d] [-g] [-h] [-m] [-p] [-v] [-z] [-D] [-d dbdir] [-r port] Infile")
	fmt.Println("   Process mixed content XML, XHTML (XDITA), and Markdown (MDITA) input.")
	fmt.Println("   Infile is a single file or directory name; no wildcards (?,*).")
	fmt.Println("          If a directory, it is processed recursively.")
	fmt.Println("   Infile may be \"-\" for Stdin: input typed (or pasted) interactively")
	fmt.Println("          is written to file ./Stdin.xml for processing")
	flag.Usage()
}

func initVars() {
	// flag.StringVar(&CA.Out.RelFilePath, "o", "", // &CA.outArg, "o", "",
	// 	"Output file name (possibly ignored, depending on command)")
	flag.StringVar(&CA.XmlCat.RelFilePath, "c", "",
		"Path to XML catalog file (do not use with \"-s\" flag)")
	flag.StringVar(&CA.DBdirPath, "d", "",
		"Directory path of/for database mmmc.db")
	flag.StringVar(&CA.XmlCatSearch.RelFilePath, "s", "",
		"Directory path to DTD schema file(s) (.dtd, .mod)")
	flag.BoolVar(&CA.DBdoImport, "m", false,
		"Import input file(s) to database")
	flag.BoolVar(&CA.FollowSymLinks, "L", true,
		"Follow symbolic links in directory recursion")
	flag.BoolVar(&CA.Pritt, "p", true,
		"Pretty-print to file with \"fmtd-\" prepended to file extension")
	flag.BoolVar(&CA.Debug, "D", false,
		"Turn on debugging")
	flag.BoolVar(&CA.Help, "h", false,
		"Show extended help message and exit")
	flag.BoolVar(&CA.GroupGenerated, "g", false,
		"Group all generated files in same-named folder \n"+
			"(e.g. ./Filnam.xml maps to ./Filenam.xml_gxml/Filenam.*)")
	flag.BoolVar(&CA.Validate, "v", true,
		"Validate input file(s)? (using xmllint) (with flag \"-c\" or \"-s\")")
	flag.BoolVar(&CA.DBdoZeroOut, "z", false,
		"Zero out the database")
	flag.IntVar(&CA.RestPort, "r", 0,
		"Run REST server on port number")
}

// checkbarf simply aborts with an error message, if a
// serious (i.e. top-level) problem has been encountered.
func checkbarf(e error, s string) {
	if e == nil {
		return
	}
	MU.SessionLogger.Printf("%s failed: %s \n", CA.AppName, e)
	fmt.Fprintf(os.Stderr, "%s failed: %s \n", CA.AppName, e)
	MU.ErrorTrace(os.Stderr, e)
	os.Exit(1)
}

// ProcessArgs wants to be a generic CLI arguments for any XML-related command.
// If this is to be so, there should be a way to selectively disable commands
// that are inappropriate for the command it is being integrated into.
func ProcessArgs(appName string, osArgs []string) (*ConfigurationArguments, error) {

	initVars()
	// Do not use logging until the invocation is sorted out.
	CA.AppName = appName
	var e error

	if !WU.IsWasm() {
		// == Figure out what CLI name we were called as ==
		osex, _ := os.Executable()
		// The call to FP.Clean(..) is needed (!!)
		println("==> Running:", FU.Tilded(FP.Clean(osex)))
		// == Locate xmllint for doing XML validations ==
		xl, e := exec.LookPath("xmllint")
		if e != nil {
			xl = "not found"
			if CA.Validate {
				println("==> Validation is not possible: xmllint cannot be found")
			}
		}
		println("==> xmllint:", xl)
	}
	// == Examine CLI invocation flags ==
	flag.Parse()
	if len(osArgs) < 2 || nil == flag.Args() || 0 == len(flag.Args()) {
		println("==> Argument parsing failed. Did not specify input file(s)?")
		myUsage()
		os.Exit(1)
	}
	if CA.Debug {
		fmt.Printf("D=> Flags: debug:%s groupGen:%s help:%s "+
			"import:%s printty:%s validate:%s zeroOutDB:%s restPort:%d \n",
			SU.Yn(CA.Debug), SU.Yn(CA.GroupGenerated), // d g h m p v z r
			SU.Yn(CA.Help), SU.Yn(CA.DBdoImport), SU.Yn(CA.Pritt),
			SU.Yn(CA.Validate), SU.Yn(CA.DBdoZeroOut), CA.RestPort)
		fmt.Println("D=> CLI tail:", flag.Args())
	}

	// Handle case where XML comes from standard input i.e. os.Stdin
	if flag.Args()[0] == "-" {
		if WU.IsWasm() {
			println("==> Trying to read from Stdin; press ^D right after a newline to end")
		} else {
			stat, e := os.Stdin.Stat()
			checkbarf(e, "Cannot Stat() Stdin")
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				println("==> Reading from Stdin; press ^D right after a newline to end")
				} else {
					println("==> Reading Stdin from a file or pipe")
				}
		}

		// bb, e := ioutil.ReadAll(os.Stdin)
		stdIn := FU.GetStringFromStdin()
		checkbarf(e, "Cannot read from Stdin")
		e = ioutil.WriteFile("Stdin.xml", []byte(stdIn), 0666)
		checkbarf(e, "Cannot write to ./Stdin.xml")
		CA.In.RelFilePath = "Stdin.xml"
	}

	// ===========================================
	//   PROCESS ARGUMENTS to get complete info
	//   about path, existence, and type
	// ===========================================

	// Process input-file(s) argument, which can be a relative filepath.
	CA.In = *FU.NewBasicPath(flag.Args()[0])

	// If the absolute path does not match the argument provided, inform the user.
	if CA.In.AbsFilePath.S() != CA.In.RelFilePath { // CA.In.ArgFilePath {
		println("==> Input:" /* FU.NiceFP */, FU.Tilded(CA.In.AbsFilePath.S()))
	}
	if CA.In.IsOkayDir() {
		println("    --> The input is a directory and will be processed recursively.")
	} else if CA.In.IsOkayFile() {
		println("    --> The input is a single file: extra info will be listed here.")
		CA.SingleFile = true
	} else {
		println("    --> The input is a type not understood.")
		return nil, fmt.Errorf("Bad type for input: " + CA.In.AbsFilePath.S())
	}

	// Process output-file(s) argument, which can be a relative filepath.
	// CA.Out.ProcessFilePathArg(CA.Out.ArgFilePath)
	CA.Out = *FU.NewBasicPath(CA.Out.RelFilePath)

	// Process database directory argument, which can be a relative filepath.
	// CA.DB.ProcessFilePathArg(CA.DBdirPath)

	// ====

	pCA := &CA
	e = pCA.ProcessDatabaseArgs()
	checkbarf(e, "Could not process DB directory argument(s)")
	e = pCA.ProcessCatalogArgs()
	checkbarf(e, "Could not process XML catalog argument(s)")
	return pCA, e
}

func (pCA *ConfigurationArguments) ProcessDatabaseArgs() error {
	var mustAccessTheDB, theDBexists bool
	var e error
	mustAccessTheDB = pCA.DBdoImport || pCA.DBdoZeroOut || pCA.DBdirPath != ""
	if !mustAccessTheDB {
		return nil
	}
	pCA.DBhandle, e = db.NewMmmcDB(pCA.DBdirPath)
	if e != nil {
		return fmt.Errorf("DB setup failure: %w", e)
	}
	theDBexists = CA.DBhandle.BasicPath.Exists
	var s = "exists"
	if !theDBexists {
		s = "does not exist"
	}
	fmt.Printf("==> DB %s: %s\n", s, pCA.DBhandle.BasicPath.AbsFilePath)

	if pCA.DBdoZeroOut {
		println("    --> Zeroing out DB")
		pCA.DBhandle.MoveCurrentToBackup()
		pCA.DBhandle.ForceEmpty()
	} else {
		pCA.DBhandle.DupeCurrentToBackup()
		pCA.DBhandle.ForceExistDBandTables()
	}
	// spew.Dump(pCA.DBhandle)
	return nil
}

func (pCA *ConfigurationArguments) ProcessCatalogArgs() error {
	var gotC, gotS bool
	gotC = ("" != CA.XmlCat.RelFilePath)
	gotS = ("" != CA.XmlCatSearch.RelFilePath)
	if !(gotC || gotS) {
		return nil
	}
	if gotC && gotS {
		return errors.New("mcfile.ConfArgs.ProcCatalArgs: cannot combine flags -c and -s")
	}
	if gotC { // -c
		// pCA.XmlCat.ProcessFilePathArg(CA.XmlCat.ArgFilePath)
		CA.XmlCat = *FU.NewBasicPath(CA.XmlCat.RelFilePath)
		if !(pCA.XmlCat.IsOkayFile() && pCA.XmlCat.Size > 0) {
			println("==> ERROR: XML catalog filepath is not file: " + pCA.XmlCat.AbsFilePath)
			return errors.New(fmt.Sprintf("mcfile.ConfArgs.ProcCatalArgs<%s:%s>",
				CA.XmlCat.RelFilePath, CA.XmlCat.AbsFilePath))
		}
		println("==> Catalog:", pCA.XmlCat.RelFilePath)
		if pCA.XmlCat.AbsFilePath.S() != pCA.XmlCat.RelFilePath {
			println("     --> i.e. ", FU.Tilded(pCA.XmlCat.AbsFilePath.S()))
		}
	}
	if gotS { // -s
		// pCA.XmlCatSearch.ProcessFilePathArg(pCA.XmlCatSearch.ArgFilePath)
		pCA.XmlCatSearch = *FU.NewBasicPath(CA.XmlCatSearch.RelFilePath)
		if !pCA.XmlCatSearch.IsOkayDir() {
			return errors.New("mcfile.ConfArgs.ProcCatalArgs: cannot open XML catalog directory: " +
				pCA.XmlCatSearch.AbsFilePath.S())
		}
	}
	var e error
	if gotS { // -s and not -c
		println("==> Schema(s):", pCA.XmlCatSearch.RelFilePath)
		// pCA.XmlCatSearch.ProcessFilePathArg(CA.XmlCatSearch.ArgFilePath)
		pCA.XmlCatSearch = *FU.NewBasicPath(CA.XmlCatSearch.RelFilePath)
		if CA.XmlCatSearch.AbsFilePath.S() != pCA.XmlCatSearch.RelFilePath {
			println("     --> i.e. ", FU.Tilded(pCA.XmlCatSearch.AbsFilePath.S()))
		}
		if !pCA.XmlCatSearch.IsOkayDir() {
			println("==> ERROR: Schema path is not a readable directory: " +
				FU.Tilded(pCA.XmlCatSearch.AbsFilePath.S()))
			return fmt.Errorf("mcfile.ConfArgs.ProcCatalArgs.abs<%s>: %w",
				pCA.XmlCatSearch.AbsFilePath, e)
		}
	}
	// println(" ")

	// ==========================
	//   PROCESS XML CATALOG(S)
	// ==========================

	// IF user asked for a single catalog file
	if gotC && !gotS {
		CA.XmlCatalog, e = gparse.NewXmlCatalogFromFile(CA.XmlCat.RelFilePath)
		if e != nil {
			println("==> ERROR: Can't find or process catalog file:", CA.XmlCat.RelFilePath)
			println("    Error was:", e.Error())
			CA.XmlCatalog = nil
			return fmt.Errorf("gxml.Confargs.NewXmlCatalogFromFile<%s>: %w", CA.XmlCat.RelFilePath, e)
		}
		if CA.XmlCatalog == nil || len(CA.XmlCatalog.XmlPublicIDs) == 0 {
			println("==> No valid entries in catalog file:", CA.XmlCat.RelFilePath)
			CA.XmlCatalog = nil
		}
		return nil
	}
	// IF user asked for a directory scan of schema files
	if gotS && !gotC {
		xmlCatalogs = make([]*gparse.XmlCatalog, 0)
		fileNameToUse := "catalog.xml"
		if CA.XmlCat.RelFilePath != "" {
			fileNameToUse = CA.XmlCat.RelFilePath
		}
		filePathToUse := FU.AbsFilePath(".")
		if CA.XmlCatSearch.RelFilePath != "" {
			filePathToUse = CA.XmlCatSearch.AbsFilePath
		}
		fileNameList, e := filePathToUse.GatherNamedFiles(fileNameToUse)
		if e != nil {
			fmt.Printf("==> No valid files named <%s> found in+under catalog search path: %s \n",
				fileNameToUse, filePathToUse)
			println("    Error was:", e.Error())
			return fmt.Errorf(
				"gxml.Confargs.GatherNamedFilesForCatalog<%s:%s>: %w", fileNameToUse, filePathToUse, e)
		}
		// For every catalog file (usually just one)
		for _, filePathToUse = range fileNameList {
			var xmlCat *gparse.XmlCatalog
			xmlCat, e = gparse.NewXmlCatalogFromFile(filePathToUse.S())
			if e != nil {
				println("==> ERROR: Can't find or process catalog file:", filePathToUse)
				println("    Error was:", e.Error())
				continue
			}
			if xmlCat == nil || len(xmlCat.XmlPublicIDs) == 0 {
				println("==> No valid entries in catalog file:", filePathToUse)
				continue
			}
			xmlCatalogs = append(xmlCatalogs, xmlCat)
		}
		switch len(xmlCatalogs) {
		case 0:
			fmt.Printf("==> ERROR: No files named <%s> found in+under <%s>:",
				fileNameToUse, filePathToUse)
			CA.XmlCatalog = nil
			return fmt.Errorf("gxml.Confargs.XmlCatalogs<%s:%s>: %w",
				fileNameToUse, filePathToUse, e)
		case 1:
			CA.XmlCatalog = xmlCatalogs[0]
		default:
			// MERGE THEM ALL
			var xmlCat *gparse.XmlCatalog
			CA.XmlCatalog = new(gparse.XmlCatalog)
			CA.XmlCatalog.XmlPublicIDs = make([]gparse.XmlPublicID, 0)
			for _, xmlCat = range xmlCatalogs {
				CA.XmlCatalog.XmlPublicIDs =
					append(CA.XmlCatalog.XmlPublicIDs, xmlCat.XmlPublicIDs...)
			}
		}
	}
	if CA.XmlCatalog == nil || CA.XmlCatalog.XmlPublicIDs == nil || len(CA.XmlCatalog.XmlPublicIDs) == 0 {
		CA.XmlCatalog = nil
		println("==> No valid catalog entries")
		return errors.New("gxml.Confargs.XmlCatalogs")
	}
	// println("==> Contents of XML catalog(s):")
	// print(CA.XmlCatalog.DString())
	fmt.Printf("==> XML catalog(s) yielded %d valid entries \n",
		len(CA.XmlCatalog.XmlPublicIDs))

	// TODO:470 If import, create batch info ?
	return nil
}
