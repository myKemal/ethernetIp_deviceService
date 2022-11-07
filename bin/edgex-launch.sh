#!/bin/bash
#
# Copyright (c) 2018
# Mainflux
#
# SPDX-License-Identifier: Apache-2.0
#

###
# Launches all EdgeX Go binaries (must be previously built).
#
# Expects that Consul is already installed and running.
#
###

DIR=$PWD
CMD=../cmd

function cleanup {
	pkill edgex-device-ethernetip-go
}

cd $CMD
exec -a edgex-device-ethernetip-go ./device-ethernetip-go &
cd $DIR


trap cleanup TERM QUIT INT

while : ; do sleep 1 ; done
