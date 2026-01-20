#!/bin/bash

while true; do
    if ! pgrep -x onboard >/dev/null; then
        onboard &
    fi
    sleep 2
done
