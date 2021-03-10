// Package experiment enables extraction of useful information from experiment objects and their formatting.
package experiment

import (
	"fmt"

	v2alpha2 "github.com/iter8-tools/etc3/api/v2alpha2"
	"gopkg.in/inf.v0"
)

// Experiment is an enhancement of v2alpha2.Experiment struct, and supports various methods used in describing an experiment.
type Experiment struct {
	v2alpha2.Experiment
}

// Started indicates if at least one iteration of the experiment has completed.
func (e *Experiment) Started() bool {
	c := e.Status.CompletedIterations
	return c != nil && *c > 0
}

// GetVersions returns the slice of version name strings. If the VersionInfo section is not present in the experiment's spec, then this slice is empty.
func (e *Experiment) GetVersions() []string {
	if e.Spec.VersionInfo == nil {
		return nil
	}
	versions := []string{e.Spec.VersionInfo.Baseline.Name}
	for _, c := range e.Spec.VersionInfo.Candidates {
		versions = append(versions, c.Name)
	}
	return versions
}

// GetMetricStr returns the metric value as a string for a given metric and a given version.
func (e *Experiment) GetMetricStr(metric string, version string) string {
	am := e.Status.Analysis.AggregatedMetrics
	if am == nil {
		return "unavailable"
	}
	if vals, ok := am.Data[metric]; ok {
		if val, ok := vals.Data[version]; ok {
			if val.Value != nil {
				z := new(inf.Dec).Round(val.Value.AsDec(), 3, inf.RoundCeil)
				return z.String()
			}
		}
	}
	return "unavailable"
}

// GetMetricStrs returns the given metric's value as a slice of strings, whose elements correspond to versions.
func (e *Experiment) GetMetricStrs(metric string) []string {
	versions := e.GetVersions()
	reqs := make([]string, len(versions))
	for i, v := range versions {
		reqs[i] = e.GetMetricStr(metric, v)
	}
	return reqs
}

// GetMetricNameAndUnits extracts the name, and if specified, units for the given metricInfo object and combines them into a string.
func GetMetricNameAndUnits(metricInfo v2alpha2.MetricInfo) string {
	r := metricInfo.Name
	if metricInfo.MetricObj.Spec.Units != nil {
		r += fmt.Sprintf(" (" + *metricInfo.MetricObj.Spec.Units + ")")
	}
	return r
}

// StringifyObjective returns a string representation of the given objective.
func StringifyObjective(objective v2alpha2.Objective) string {
	r := ""
	if objective.LowerLimit != nil {
		z := new(inf.Dec).Round(objective.LowerLimit.AsDec(), 3, inf.RoundCeil)
		r += z.String() + " <= "
	}
	r += objective.Metric
	if objective.UpperLimit != nil {
		z := new(inf.Dec).Round(objective.UpperLimit.AsDec(), 3, inf.RoundCeil)
		r += " <= " + z.String()
	}
	return r
}

// GetSatisfyStr returns a true/false/unavailable valued string denotating if a version satisfies the objective.
func (e *Experiment) GetSatisfyStr(objectiveIndex int, version string) string {
	ana := e.Status.Analysis
	if ana == nil {
		return "unavailable"
	}
	va := ana.VersionAssessments
	if va == nil {
		return "unavailable"
	}
	if vals, ok := va.Data[version]; ok {
		if len(vals) > objectiveIndex {
			return fmt.Sprintf("%v", vals[objectiveIndex])
		}
	}
	return "unavailable"
}

// GetSatisfyStrs returns a slice of true/false/unavailable valued strings for an objective denoting if it is satisfied by versions.
func (e *Experiment) GetSatisfyStrs(objectiveIndex int) []string {
	versions := e.GetVersions()
	sat := make([]string, len(versions))
	for i, v := range versions {
		sat[i] = e.GetSatisfyStr(objectiveIndex, v)
	}
	return sat
}
