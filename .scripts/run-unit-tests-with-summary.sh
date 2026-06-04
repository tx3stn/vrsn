#! /bin/sh

# Script to run unit tests in CI and gracefully handle the status to display
# all of the info you need to debug and output the results to the
# GITHUB_STEP_SUMMARY environment variable so the test summary is always
# displayed if things pass or fail.

TEST_OUTPUT=$(CGO_ENABLED=1 go test ./... -race -cover)
EXIT_CODE="$?"

if [ "$EXIT_CODE" = 0 ]; then
	echo "# 🎉 unit tests passed 🎉"
	echo ""
	SUMMARY=$(echo "$TEST_OUTPUT" | ./.scripts/parse-unit-tests.sh)
	echo "$SUMMARY"
	exit 0
fi

echo "# 🚩 some tests failed 🚩"
echo "\`\`\`"
echo "$TEST_OUTPUT"
echo "\`\`\`"
exit 1
