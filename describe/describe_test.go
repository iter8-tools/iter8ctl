package describe

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/iter8-tools/iter8ctl/utils"
	"github.com/stretchr/testify/assert"
)

/* Tests */

func TestBuilder(t *testing.T) {
	a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
	d := Builder(&a, &b, &c)
	d.Usage()
	assert.Greater(t, c.Len(), 0)
	assert.NoError(t, d.Error())
}

func TestInvalidArguments(t *testing.T) {
	for _, args := range [][]string{
		{"-name", "helloworld"},
		{},
		{"-f"},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		d := Builder(&a, &b, &c)
		d.ParseFlags(args)
		assert.Error(t, d.Error())
		e := d.Error()
		d.GetExperiment()
		assert.Error(t, d.Error())
		assert.Equal(t, e, d.Error())
	}
}

func TestFileInputGood(t *testing.T) {
	for _, args := range [][]string{
		{"-f", utils.CompletePath("../", "testdata/experiment1.yaml")},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		d := Builder(&a, &b, &c)
		d.ParseFlags(args).GetExperiment()
		assert.NoError(t, d.Error())
	}
}

func TestStdinGood(t *testing.T) {
	for _, args := range [][]string{
		{"-f", "-"},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		data, _ := ioutil.ReadFile(utils.CompletePath("../", "testdata/experiment1.yaml"))
		a.Write(data)
		d := Builder(&a, &b, &c)
		d.ParseFlags(args).GetExperiment()
		assert.NoError(t, d.Error())
	}
}

func TestFileDoesNotExist(t *testing.T) {
	for _, args := range [][]string{
		{"-f", utils.CompletePath("../", "testdata/nonexistant.yaml")},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		d := Builder(&a, &b, &c)
		d.ParseFlags(args).GetExperiment()
		assert.Error(t, d.Error())
	}
}

func TestStdinBadYAML(t *testing.T) {
	for _, args := range [][]string{
		{"-f", "-", "abc 123 xyz"},
		{"-f", "-", "abc:\nabc:1"},
	} {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		a.Write([]byte(args[2]))
		d := Builder(&a, &b, &c)
		d.ParseFlags(args[:2]).GetExperiment()
		assert.Error(t, d.Error())
	}
}

func TestPrintProgress(t *testing.T) {
	for i := 1; i <= 12; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		d := Builder(&a, &b, &c)
		d.ParseFlags([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().printProgress()
		assert.NoError(t, d.Error())
	}
}

func TestPrintWinnerAssessment(t *testing.T) {
	for i := 1; i <= 12; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		d := Builder(&a, &b, &c)
		d.ParseFlags([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().printWinnerAssessment()
		assert.NoError(t, d.Error())
	}
}

func TestPrintObjectiveAssessment(t *testing.T) {
	for i := 1; i <= 12; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		d := Builder(&a, &b, &c)
		d.ParseFlags([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().printObjectiveAssessment()
		assert.NoError(t, d.Error())
	}
}

func TestPrintVersionAssessment(t *testing.T) {
	for i := 1; i <= 12; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		d := Builder(&a, &b, &c)
		d.ParseFlags([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().printVersionAssessment()
		assert.NoError(t, d.Error())
	}
}

func TestPrintMetrics(t *testing.T) {
	for i := 1; i <= 12; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		d := Builder(&a, &b, &c)
		d.ParseFlags([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().printMetrics()
		assert.NoError(t, d.Error())
	}
}

func TestPrintRewardAssessments(t *testing.T) {
	for i := 1; i <= 12; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		d := Builder(&a, &b, &c)
		d.ParseFlags([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().printRewardAssessment()
		assert.NoError(t, d.Error())
	}
}

func TestPrintAnalysis(t *testing.T) {
	for i := 1; i <= 12; i++ {
		a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
		d := Builder(&a, &b, &c)
		d.ParseFlags([]string{"-f", utils.CompletePath("../", fmt.Sprintf("testdata/experiment%v.yaml", i))}).GetExperiment().PrintAnalysis()
		assert.NoError(t, d.Error())
	}
}

/* Examples */

func ExampleBuilder() {
	d := Builder(os.Stdin, os.Stdout, os.Stderr)
	d.ParseFlags([]string{"-f", "path-to-my-experiment.yaml"})
}

func ExampleBuilder_bytebuffers() {
	a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
	d := Builder(&a, &b, &c)
	// The following will print the usage message in d's stderr, i.e., byte buffer c.
	d.Usage()
}

func ExampleCmd_Error() {
	d := Builder(os.Stdin, os.Stdout, os.Stderr)
	// "-g" is an invalid flag which will cause ParseFlags invocation to generate an error.
	d.ParseFlags([]string{"-g", "golly"})
	// This will print the error to d.stderr (= os.Stderr)
	fmt.Fprintln(d.stderr, d.Error())
}

func ExampleCmd_Error_bytebuffers() {
	a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
	d := Builder(&a, &b, &c)
	// "-g" is an invalid flag which will cause ParseFlags invocation to generate an error.
	d.ParseFlags([]string{"-g", "golly"})
	// The following will print the error message in d's stderr, i.e., byte buffer c.
	fmt.Fprintln(d.stderr, d.Error())
}
func ExampleCmd_GetExperiment() {
	d := Builder(os.Stdin, os.Stdout, os.Stderr)
	// Invalid experiment input will cause GetExperiment to generate an error
	d.ParseFlags([]string{"-f", "path-to-my-experiment.yaml"}).GetExperiment()
}

func ExampleCmd_ParseFlags() {
	d := Builder(os.Stdin, os.Stdout, os.Stderr)
	// "-f" is the only supported flag
	d.ParseFlags([]string{"-f", "path-to-my-experiment.yaml"})
}

func ExampleCmd_ParseFlags_invalid() {
	d := Builder(os.Stdin, os.Stdout, os.Stderr)
	// Invalid flags will cause ParseFlags to generate an error.
	d.ParseFlags([]string{"-g", "golly"})
}

func ExampleCmd_PrintAnalysis() {
	d := Builder(os.Stdin, os.Stdout, os.Stderr)
	d.ParseFlags([]string{"-f", "path-to-my-experiment.yaml"}).
		GetExperiment().
		PrintAnalysis()
}

func ExampleCmd_PrintAnalysis_bytebuffers() {
	a, b, c := bytes.Buffer{}, bytes.Buffer{}, bytes.Buffer{}
	d := Builder(&a, &b, &c)
	// PrintAnalysis call below will print to d.stdout, i.e., to byte buffer b.
	d.ParseFlags([]string{"-f", "path-to-my-experiment.yaml"}).
		GetExperiment().
		PrintAnalysis()
}
