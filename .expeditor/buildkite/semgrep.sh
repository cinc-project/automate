#!/bin/bash
set -eou pipefail

# usage: semgrep.sh $BUILDKITE_BRANCH $BUILDKITE_BUILD_URL $BUILDKITE_PULL_REQUEST_BASE_BRANCH

echo "running in $(pwd)"
set -x
export SEMGREP_BRANCH=$1
export SEMGREP_REPO_NAME=chef/automate
# Activates links to buildkite builds in slack notification
export SEMGREP_JOB_URL=$2
# Activates links to git commits in slack notification
export SEMGREP_REPO_URL=https://github.com/chef/automate
BASELINE=$3
MERGE_BASE=$(git merge-base ${BASELINE:-master} HEAD)
set +x
# \[ "$BUILDKITE_BRANCH" != "master" \] && \[ -n "\$BASELINE" \] && echo "PR build on '$BUILDKITE_BRANCH' branch; base is \$MERGE_BASE (merge-base of '\$BASELINE')" && BASELINE=\$MERGE_BASE
# \[ "$BUILDKITE_BRANCH" != "master" \] && \[ -z "\$BASELINE" \] && echo "manual build on '$BUILDKITE_BRANCH' branch; using merge-base of master as base (\$MERGE_BASE)" && BASELINE=\$MERGE_BASE
# \[ "$BUILDKITE_BRANCH" == "master" \] && echo "build on master; using 'master~1' as base branch" && BASELINE=master~1
# python -m semgrep_agent --publish-token "\$SEMGREP_TOKEN" --publish-deployment \$SEMGREP_ID --baseline-ref \$BASELINE