#!/usr/bin/env bash

# Workflow environment variables
# These variables create an Alfred-like environment

root="$( git rev-parse --show-toplevel )"
testdir="${root}/testenv"

# Absolute bare-minimum for AwGo to function...
export alfred_workflow_bundleid="com.github.dongc.prompts"
export alfred_workflow_data="${testdir}/data"
export alfred_workflow_cache="${testdir}/cache"
echo $alfred_workflow_bundleid
echo 

test -f "$HOME/Library/Preferences/com.runningwithcrayons.Alfred.plist" || {
	export alfred_version="4.6.6"
}

# Expected by ExampleNew
export alfred_workflow_version="1.0.0"
export alfred_workflow_name="alfred-dndbeyond-monster-workflow"

# Prevent random ID from being generated
export AW_SESSION_ID="local-session"