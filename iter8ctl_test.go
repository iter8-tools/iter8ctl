package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/iter8-tools/iter8ctl/utils"
	"github.com/stretchr/testify/assert"
)

type test struct {
	name           string   // name of this test
	flags          []string // flags supplied to .iter8ctl command
	outputFilename string   // relative to testdata
	errorFilename  string   // relative to testdata
}

var tests = []test{
	{name: "no-flags", flags: []string{}, outputFilename: "", errorFilename: "error-no-flags.txt"},
	{name: "invalid-subcommand", flags: []string{"invalid"}, outputFilename: "", errorFilename: "error-invalid-subcommand.txt"},
	{name: "undefined-flag", flags: []string{"describe", "-name", "helloworld"}, outputFilename: "", errorFilename: "error-undefined-flag.txt"},
	{name: "experiment1", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment1.yaml")}, outputFilename: "experiment1.out", errorFilename: ""},
	{name: "experiment2", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment2.yaml")}, outputFilename: "experiment2.out", errorFilename: ""},
	{name: "experiment3", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment3.yaml")}, outputFilename: "experiment3.out", errorFilename: ""},
	{name: "experiment4", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment4.yaml")}, outputFilename: "experiment4.out", errorFilename: ""},
	{name: "experiment5", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment5.yaml")}, outputFilename: "experiment5.out", errorFilename: ""},
	{name: "experiment6", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment6.yaml")}, outputFilename: "experiment6.out", errorFilename: ""},
	{name: "experiment7", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment7.yaml")}, outputFilename: "experiment7.out", errorFilename: ""},
	{name: "experiment8", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment8.yaml")}, outputFilename: "experiment8.out", errorFilename: ""},
	{name: "experiment9", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment9.yaml")}, outputFilename: "experiment9.out", errorFilename: ""},
	{name: "experiment11", flags: []string{"describe", "-f", utils.CompletePath("testdata", "experiment11.yaml")}, outputFilename: "experiment11.out", errorFilename: ""},
}

/* CLI tests */

func TestCLI(t *testing.T) {
	// All subtests within this test rely in ./iter8ctl. So go build.
	exec.Command("go", "build", utils.CompletePath("", "")).Run()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cmd := exec.Command("./iter8ctl", tc.flags...)
			b1 := &strings.Builder{}
			cmd.Stderr = b1
			b2 := &strings.Builder{}
			cmd.Stdout = b2

			cmd.Run()

			if tc.errorFilename != "" {
				ef := utils.CompletePath("testdata", tc.errorFilename)
				t.Log("Reading error file: ", ef)
				b3, err := ioutil.ReadFile(utils.CompletePath("testdata", tc.errorFilename))
				if err != nil {
					t.Fatal("Unable to read contents of error file: ", ef)
				}
				assert.Equal(t, string(b3), b1.String())
			}

			if tc.outputFilename != "" {
				of := utils.CompletePath("testdata", tc.outputFilename)
				t.Log("Reading output file: ", of)
				b4, err := ioutil.ReadFile(utils.CompletePath("testdata", tc.outputFilename))
				if err != nil {
					t.Fatal("Unable to read contents of output file: ", of)
				}
				assert.Equal(t, string(b4), b2.String())
			}
		})
	}
}

/* main() tests -- identical to CLI tests, except, now we will invoke main */

func TestMain(t *testing.T) {
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			os.Args = append([]string{"./iter8ctl"}, tc.flags...)
			b1 := &strings.Builder{}
			stderr = b1
			b2 := &strings.Builder{}
			stdout = b2

			main()

			if tc.errorFilename != "" {
				ef := utils.CompletePath("testdata", tc.errorFilename)
				b3, err := ioutil.ReadFile(ef)
				if err != nil {
					t.Fatal("Unable to read contents of error file: ", ef)
				}
				assert.Equal(t, string(b3), b1.String())
			}

			if tc.outputFilename != "" {
				of := utils.CompletePath("testdata", tc.outputFilename)
				b4, err := ioutil.ReadFile(of)
				if err != nil {
					t.Fatal("Unable to read contents of output file: ", of)
				}
				assert.Equal(t, string(b4), b2.String())
			}
		})
	}
}
