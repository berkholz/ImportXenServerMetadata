package main

// ################## imports
import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// ################## global variables

const appVersion string = "0.1"

// filename with xml meta data information
var importFile string = ""

// xe binary which is used to import the metadata
var xeBinary string = ""

// root element of our vms
type vms struct {
	XMLName xml.Name `xml:"vms"`
	Vms     []vm     `xml:"vm"`
}

// a vm
type vm struct {
	XMLName   xml.Name  `xml:"vm"`
	NameLabel string    `xml:"name,attr"`
	UUID      string    `xml:"uuid,attr"`
	Parents   parents   `xml:"parents"`
	VBDs      vbds      `xml:"vbds"`
	Snapshots snapshots `xml:"snapshots"`
}

// list of parents of a vm
type parents struct {
	XMLName xml.Name `xml:"parents"`
	Parents []parent `xml:"parent"`
}

// parent of a vm, part of the list
type parent struct {
	XMLName    xml.Name `xml:"parent"`
	UUID       string   `xml:"uuid,attr"`
	Selfparent bool     `xml:"selfparent,attr"`
}

// list of vbds of a vm
type vbds struct {
	XMLName xml.Name `xml:"vbds"`
	VBDs    []vbd    `xml:"vbd"`
}

// vbd of a vm, part of the list
type vbd struct {
	XMLName      xml.Name `xml:"vbd"`
	UUID         string   `xml:"uuid,attr"`
	VBDType      string   `xml:"type,attr"`
	VdiNameLabel string   `xml:"vdi-name-label,attr"`
}

// list of snapshots of a vm
type snapshots struct {
	XMLName   xml.Name   `xml:"snapshots"`
	Snapshots []snapshot `xml:"snapshot"`
}

// snapshot of a vm, part of the list
type snapshot struct {
	XMLName         xml.Name `xml:"snapshot"`
	UUID            string   `xml:"uuid,atrr"`
	NameLable       string   `xml:"name-lable,atrr"`
	NameDescription string   `xml:"name-description,atrr"`
	IsVmssSnapshot  string   `xml:"is-vmss-snapshot,atrr"`
}

func parseCommandOptions() {
	flag.StringVar(&xeBinary, "xebinary", xeBinary, "Absolute path to xe binary including executable.")

	flag.StringVar(&importFile, "infile", importFile, "Filename with meta data to import.")

	// The flag package provides a default help printer via -h switch
	var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

	flag.Parse() // Scan the arguments list

	// check if flag -h is given, print out version end exit.
	if *versionFlag {
		fmt.Println("Version:", appVersion)
		os.Exit(0)
	}

	// check xe binary
	if len(xeBinary) == 0 {
		xeBinary = "/usr/bin/xe"
		log.Println("No path to xe binary given on command line using default: " + xeBinary)
	}
	path, err := exec.LookPath(xeBinary)
	if err != nil {
		log.Println("Could not find xe binary. Exiting...")
	} else {
		fmt.Printf("xe binary found in %v", path)
	}

	// check if an import file name is specified
	if len(importFile) == 0 {
		importFile = "vms.export.xml"
		fmt.Println("No export file given, using default: " + importFile)
	}
}

func main() {
	// ################## variable definitions

	parseCommandOptions()

	// we initialize our Users array
	var virtualMachines vms

	// read in the import file
	xmlInputFile, err := os.Open(importFile)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully Opened " + importFile)
	}

	defer xmlInputFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlInputFile)

	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &virtualMachines)

	// fmt.Printf("--> VMs %+v", virtualMachines)

	for index, element := range virtualMachines.Vms {
		fmt.Printf("%v : ", index)
		// fmt.Println(element.VBDs)
		fmt.Printf("%v\n", element.UUID)
	}
}
