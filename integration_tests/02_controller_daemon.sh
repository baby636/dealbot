#!/bin/bash

my_dir="$(dirname "$0")"
source "$my_dir/header.sh"

export DEALBOT_MINER_ADDRESS=t01000

function stop_dealbot_and_devnet {
	kill -9 $CONTROLLER_PID $DAEMON_PID $CONTROLLER_TAIL_PID $DAEMON_TAIL_PID

	# We replace devnet's EXIT trap, so run it too.
	stop_devnet
}
trap stop_dealbot_and_devnet EXIT

CONTROLLER_LISTEN="localhost:8085"
CONTROLLER_ENDPOINT="http://${CONTROLLER_LISTEN}"

dealbot controller --listen="${CONTROLLER_LISTEN}" &>dealbot-controller.log &
CONTROLLER_PID=$!

tail -f dealbot-controller.log | sed 's/^/dealbot-controller.log: /' &
CONTROLLER_TAIL_PID=$!

# Give it half a second to start.
sleep 0.5
if ! kill -0 $CONTROLLER_PID; then
	tail -n 50 dealbot-controller.log
	unset CONTROLLER_PID # no need to kill it
	exit 1
fi

dealbot daemon --endpoint=$CONTROLLER_ENDPOINT &>dealbot-daemon.log  &
DAEMON_PID=$!

# Give it half a second to start.
sleep 0.5
if ! kill -0 $DAEMON_PID; then
	tail -n 50 dealbot-daemon.log
	unset DAEMON_PID # no need to kill it
	exit 1
fi

tail -f dealbot-daemon.log | sed 's/^/dealbot-daemon.log: /' &
DAEMON_TAIL_PID=$!

curl --header "Content-Type: application/json" \
	--request POST \
	--data '{"Miner":"t01000","MaxPriceAttoFIL":100000000000000000,"Size":1024,"StartOffset":0,"FastRetrieval":true,"Verified":false}' \
	"$CONTROLLER_ENDPOINT/tasks/storage"

# TODO: Also test a retrieval task.

# TODO: Once the daemon posts all updates back to the controller, use the
# controller logs alone.

# TODO: Track the status logs more intelligently. Probably too complex for shell.

# TODO: Also test multiple daemons, and e.g. submitting many trivial tasks and
# seeing that every daemon picks up at least one of them.

# On average, the storage deal takes about four minutes.
# Use a timeout of ten minutes, just in case.
for ((i = 0; ; i++)); do
	if grep -q StorageDealActive dealbot-daemon.log; then
		echo "Storage deal is active!"
		break
	fi
	if ((i > 10*60)); then
		# The logs are already being printed out.
		echo "Timed out waiting for storage deal to be active."
		exit 0
	fi
	sleep 1
done

if ! grep -q 'INFO.*controller.*storage task.*"miner": "t01000",.*"size": 1024,.*"status": 2,' dealbot-controller.log; then
	echo "The controller doesn't seem to be logging tasks correctly."
	exit 1
fi
