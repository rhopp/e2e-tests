#!/bin/bash

# exit immediately when a command fails
set -e
# only exit with zero if all commands of the pipeline exit successfully
set -o pipefail
# error on unset variables
set -u

export WORKSPACE=$(dirname $(dirname $(readlink -f "$0")));

# Available openshift ci environments https://docs.ci.openshift.org/docs/architecture/step-registry/#available-environment-variables
export ARTIFACTS_DIR=${ARTIFACT_DIR:-"/tmp/appstudio"}

function executeE2ETests() {
    ginkgo -p "${WORKSPACE}"/cmd --junit-report="${ARTIFACTS_DIR}"/e2e-report.xml --progress --v --label-filter='e2e-demo' -- --config-suites="${WORKSPACE}"/tests/e2e-demos/config/default.yaml
}

go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo@latest
go mod tidy -compat=1.17
go mod vendor
ginkgo version

# Initiate openshift ci users
export KUBECONFIG_TEST="/tmp/kubeconfig"
/bin/bash "$WORKSPACE"/scripts/provision-openshift-user.sh

export KUBECONFIG="${KUBECONFIG_TEST}"

/bin/bash "$WORKSPACE"/scripts/install-appstudio-e2e-mode.sh install

executeE2ETests
