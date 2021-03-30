// Package describe implements the `iter8ctl describe` subcommand.
package describe

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/iter8-tools/etc3/api/v2alpha2"
	"github.com/iter8-tools/iter8ctl/experiment"
	"github.com/olekukonko/tablewriter"
)

// Cmd struct contains fields that store flags and intermediate results associated with an invocation of 'iter8ctl describe' subcommand.
type Cmd struct {
	flagSet        *flag.FlagSet
	experimentPath string
	experiment     *experiment.Experiment
	description    strings.Builder
	err            error
	stdin          io.Reader
	stdout         io.Writer
	stderr         io.Writer
	// Usage is a function that is invoked when execution of any `Cmd` method results in an error.
	// The typical behavior of `Cmd` after an error is as follows: Usage() prints an error message to stderr, subsequent Cmd methods turn into no-ops, and the program exits.
	// Note that Usage() is a field and not a method. You can supply your own implementation of Usage() while constructing `Cmd`.
	Usage func()
}

// Builder returns an initialized Cmd struct pointer.
// Builder enables the builder design pattern along with method chaining.
func Builder(stdin io.Reader, stdout io.Writer, stderr io.Writer) *Cmd {
	// ContinueOnError ensures flagSet does not ExitOnError
	var flagSet = flag.NewFlagSet("describe", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	var d = &Cmd{
		flagSet:        flagSet,
		experimentPath: "",
		experiment:     &experiment.Experiment{},
		description:    strings.Builder{},
		err:            nil,
		stdin:          stdin,
		stdout:         stdout,
		stderr:         stderr,
		Usage:          flagSet.Usage,
	}

	//setup flagSet
	const (
		defaultExperimentPath = ""
		usage                 = "absolute path to experiment yaml file, or - for console input (stdin)"
	)
	flagSet.StringVar(&d.experimentPath, "f", defaultExperimentPath, usage)

	return d
}

// Error returns any error generated during the invocation of Cmd methods, or nil if there are no errors.
func (d *Cmd) Error() error {
	return d.err
}

// ParseFlags parses the flags supplied to Cmd. The returned Cmd struct contains the parsed result. If invalid flags are supplied, ParseFlags generates an error.
func (d *Cmd) ParseFlags(args []string) *Cmd {
	if d.err != nil {
		return d
	}
	d.err = d.flagSet.Parse(args)
	if d.err == nil {
		if d.flagSet.NFlag() == 0 {
			d.err = errors.New("Missing experiment")
			fmt.Fprintln(d.stderr, d.err)
			d.Usage()
		}
	}
	return d
}

// GetExperiment populates the Cmd struct with an experiment.
// The experiment may come from an input file when `iter8ctl describe` subcommand is invoked with the "-f experiment-file-path.yaml" flag.
// The experiment may also come from console input when `iter8ctl describe` subcommand is invoked with the "-f -" flag.
// The experiment input needs to be a valid iter8 experiment YAML. Otherwise, GetExperiment will generate an error.
func (d *Cmd) GetExperiment() *Cmd {
	if d.err != nil {
		return d
	}
	var expBytes []byte
	if d.experimentPath == "-" {
		expBytes, d.err = ioutil.ReadAll(d.stdin)
	} else {
		expBytes, d.err = ioutil.ReadFile(d.experimentPath)
	}
	if d.err != nil {
		d.err = errors.New("Error reading experiment YAML input")
		fmt.Fprintln(d.stderr, d.err)
		return d
	}
	if len(expBytes) == 0 {
		d.err = errors.New("Error reading experiment YAML input... zero bytes read")
		fmt.Fprintln(d.stderr, d.err)
		return d
	}
	d.err = yaml.Unmarshal(expBytes, d.experiment)
	if d.err != nil {
		d.err = errors.New("unmarshal error... this could be due to invalid experiment YAML input")
		fmt.Fprintln(d.stderr, d.err)
		return d
	}

	return d
}

// printProgress prints name, namespace, and target of the experiment and the number of completed iterations into d's description buffer.
func (d *Cmd) printProgress() *Cmd {
	if d.err != nil {
		return d
	}
	d.description.WriteString("\n****** Overview ******\n")
	d.description.WriteString("Experiment name: " + d.experiment.Name + "\n")
	d.description.WriteString("Experiment namespace: " + d.experiment.Namespace + "\n")
	d.description.WriteString("Target: " + d.experiment.Spec.Target + "\n")
	d.description.WriteString(fmt.Sprintf("Testing pattern: %v\n", d.experiment.Spec.Strategy.TestingPattern))
	if d.experiment.Spec.Strategy.DeploymentPattern != nil {
		d.description.WriteString(fmt.Sprintf("Deployment pattern: %v\n", *d.experiment.Spec.Strategy.DeploymentPattern))
	}

	d.description.WriteString("\n****** Progress Summary ******\n")
	sta := d.experiment.Status
	if sta.Stage != nil {
		d.description.WriteString(fmt.Sprintf("Experiment stage: %s\n", *sta.Stage))
	}
	if sta.CompletedIterations == nil || *sta.CompletedIterations == 0 {
		d.description.WriteString("Number of completed iterations: 0\n")
	} else {
		d.description.WriteString(fmt.Sprintf("Number of completed iterations: %v\n", *sta.CompletedIterations))
	}
	return d
}

// printWinnerAssessment prints the winning version in the experiment into d's description buffer.
// If winner assessment is unavailable for the underlying experiment, this method will indicate likewise.
func (d *Cmd) printWinnerAssessment() *Cmd {
	if d.err != nil {
		return d
	}
	if a := d.experiment.Status.Analysis; a != nil {
		if w := a.WinnerAssessment; w != nil {
			d.description.WriteString("\n****** Winner Assessment ******\n")
			var explanation string = ""
			switch d.experiment.Spec.Strategy.TestingPattern {
			case v2alpha2.TestingPatternCanary:
				explanation = "> If the candidate version satisfies the experiment objectives, then it is the winner.\n> Otherwise, if the baseline version satisfies the experiment objectives, it is the winner.\n> Otherwise, there is no winner.\n"
			case v2alpha2.TestingPatternConformance:
				explanation = "> If the version being validated; i.e., the baseline version, satisfies the experiment objectives, it is the winner.\n> Otherwise, there is no winner.\n"
			default:
				explanation = ""
			}
			d.description.WriteString(explanation)
			if d.experiment.Spec.Strategy.TestingPattern != v2alpha2.TestingPatternConformance && d.experiment.Spec.VersionInfo != nil {
				versions := []string{d.experiment.Spec.VersionInfo.Baseline.Name}
				for i := 0; i < len(d.experiment.Spec.VersionInfo.Candidates); i++ {
					versions = append(versions, d.experiment.Spec.VersionInfo.Candidates[i].Name)
				}
				d.description.WriteString(fmt.Sprintf("App versions in this experiment: %s\n", versions))
			}
			if w.Data.WinnerFound {
				d.description.WriteString(fmt.Sprintf("Winning version: %s\n", *w.Data.Winner))
			} else {
				d.description.WriteString("Winning version: not found\n")
			}

			if d.experiment.Spec.Strategy.TestingPattern != v2alpha2.TestingPatternConformance &&
				d.experiment.Status.VersionRecommendedForPromotion != nil {
				d.description.WriteString(fmt.Sprintf("Version recommended for promotion: %s\n", *d.experiment.Status.VersionRecommendedForPromotion))
			}
		}
	}
	return d
}

// printObjectiveAssessment prints a matrix of boolean values into d's description buffer.
// Rows correspond to experiment objectives, columns correspond to versions, and entry [i, j] indicates if objective i is satisfied by version j.
// Objective assessments are printed in the same sequence as in the experiment's spec.criteria.objectives section.
// If objective assessments are unavailable for the underlying experiment, this method will indicate likewise.
func (d *Cmd) printObjectiveAssessment() *Cmd {
	if d.err != nil {
		return d
	}
	if a := d.experiment.Status.Analysis; a != nil {
		if v := a.VersionAssessments; v != nil {
			d.description.WriteString("\n****** Objective Assessment ******\n")
			d.description.WriteString("> Identifies whether or not the experiment objectives are satisfied by the most recently observed metrics values for each version.\n")
			table := tablewriter.NewWriter(&d.description)
			table.SetRowLine(true)
			versions := d.experiment.GetVersions()
			table.SetHeader(append([]string{"Objective"}, versions...))
			for i, objective := range d.experiment.Spec.Criteria.Objectives {
				row := []string{experiment.StringifyObjective(objective)}
				table.Append(append(row, d.experiment.GetSatisfyStrs(i)...))
			}
			table.Render()
		}
	}
	return d
}

// printVersionAssessment prints how each version is performing with respect to experiment criteria into d's description buffer. This method invokes printObjectiveAssessment under the covers.
func (d *Cmd) printVersionAssessment() *Cmd {
	if d.err != nil {
		return d
	}
	if c := d.experiment.Spec.Criteria; c != nil && len(c.Objectives) > 0 {
		d.printObjectiveAssessment()
	}
	return d
}

// printMetrics prints a matrix of (decimal) metric values into d's description buffer.
// Rows correspond to experiment metrics, columns correspond to versions, and entry [i, j] indicates the value of metric i for version j.
// Metrics are printed in the same sequence as in the experiment's status.metrics section.
// If metrics are unavailable for the underlying experiment, this method will indicate likewise.
func (d *Cmd) printMetrics() *Cmd {
	if d.err != nil {
		return d
	}
	if a := d.experiment.Status.Analysis; a != nil {
		if v := a.AggregatedMetrics; v != nil {
			d.description.WriteString("\n****** Metrics Assessment ******\n")
			d.description.WriteString("> Most recently read values of experiment metrics for each version.\n")
			table := tablewriter.NewWriter(&d.description)
			table.SetRowLine(true)
			versions := d.experiment.GetVersions()
			table.SetHeader(append([]string{"Metric"}, versions...))
			for _, metricInfo := range d.experiment.Status.Metrics {
				row := []string{experiment.GetMetricNameAndUnits(metricInfo)}
				table.Append(append(row, d.experiment.GetMetricStrs(metricInfo.Name)...))
			}
			table.Render()
		}
	}
	return d
}

// PrintAnalysis prints the progress of the iter8 experiment, winner assessment, version assessment, and metrics.
func (d *Cmd) PrintAnalysis() *Cmd {
	if d.err != nil {
		return d
	}
	d.printProgress()
	if d.experiment.Started() {
		d.printWinnerAssessment().printVersionAssessment().printMetrics()
	}
	if d.err == nil {
		fmt.Fprintln(d.stdout, d.description.String())
	}
	return d
}
