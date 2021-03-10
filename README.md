[![Go Report Card](https://goreportcard.com/badge/github.com/iter8-tools/iter8ctl)](https://goreportcard.com/report/github.com/iter8-tools/iter8ctl)
[![Coverage](https://codecov.io/gh/iter8-tools/iter8ctl/branch/main/graphs/badge.svg?branch=main)](https://codecov.io/gh/iter8-tools/iter8ctl)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/iter8-tools/iter8ctl.svg)](https://pkg.go.dev/github.com/iter8-tools/iter8ctl)
# Iter8ctl
Iter8 command line utility for service operators to understand and diagnose their iter8 experiments.

Iter8ctl can be used with [iter8-kfserving](https://github.com/iter8-tools/iter8-kfserving) experiments.

# Installation
```
GOBIN=/usr/local/bin go get github.com/iter8-tools/iter8ctl
```
The above command installs `iter8ctl` under the `/usr/local/bin` directory. To install under a different directory, change the value of `GOBIN` above.

## Removal
```
rm <path-to-install-directory>/iter8ctl
```

# Usage

## Example 1
Describe an iter8 Experiment resource object present in your Kubernetes cluster.
```shell
kubectl get experiment sklearn-iris-experiment-1 -n kfserving-test -o yaml > experiment.yaml
iter8ctl describe -f experiment.yaml
```

## Example 2
Supply experiment YAML using console input.
```shell
kubectl get experiment sklearn-iris-experiment-1 -n kfserving-test -o yaml > experiment.yaml
cat experiment.yaml | iter8ctl describe -f -
```

## Example 3
Periodically fetch an iter8 Experiment resource object present in your Kubernetes cluster and describe it. You can change the frequency by adjusting the sleep interval below.
```shell
while clear; do
    kubectl get experiment sklearn-iris-experiment-1 -n kfserving-test -o yaml | iter8ctl describe -f -
    sleep 10.0
done
```

# Sample Output
The following is the output of executing `iter8ctl describe -f testdata/experiment10.yaml`; the `testdata` folder is part of the [iter8ctl GitHub repo](https://github.com/iter8-tools/iter8ctl) and contains sample experiments used in tests.

```shell
$ ./iter8ctl describe -f testdata/experiment10.yaml
****** Overview ******
Experiment name: experiment-1
Experiment namespace: knative-test
Target: knative-test/sample-application
Testing pattern: Canary
Deployment pattern: Progressive

****** Progress Summary ******
Experiment stage: Completed
Number of completed iterations: 8

****** Winner Assessment ******
App versions in this experiment: [sample-application-v1 sample-application-v2]
Winning version: sample-application-v2
Version recommended for promotion: sample-application-v2

****** Objective Assessment ******
+--------------------------+-----------------------+-----------------------+
|        OBJECTIVE         | SAMPLE-APPLICATION-V1 | SAMPLE-APPLICATION-V2 |
+--------------------------+-----------------------+-----------------------+
| mean-latency <= 2000.000 | true                  | true                  |
+--------------------------+-----------------------+-----------------------+
| error-rate <= 0.010      | true                  | true                  |
+--------------------------+-----------------------+-----------------------+

****** Metrics Assessment ******
+-----------------------------+-----------------------+-----------------------+
|           METRIC            | SAMPLE-APPLICATION-V1 | SAMPLE-APPLICATION-V2 |
+-----------------------------+-----------------------+-----------------------+
| request-count               |              1022.565 |               514.445 |
+-----------------------------+-----------------------+-----------------------+
| mean-latency (milliseconds) |                 5.881 |                 4.702 |
+-----------------------------+-----------------------+-----------------------+
| error-rate                  |                 0.000 |                 0.000 |
+-----------------------------+-----------------------+-----------------------+
```

# Contributing

Documentation and code PRs are welcome. When contributing to this repository, please first discuss the change you wish to make using [Issues](https://github.com/iter8-tools/iter8ctl/issues), [Discussion](https://github.com/iter8-tools/iter8ctl/discussions), or [Slack](https://join.slack.com/t/iter8-tools/shared_invite/enQtODU0NTczMTQ5NDU4LTJmNGE1OTBhOWI4NzllZGE0ZjdhM2M3MzJlMjcxYjliMTJlM2YxMzQ4OWQ5NGViYTM2MTU4MWRkZTgxNzZiMzg).