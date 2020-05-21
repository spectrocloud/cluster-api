#!/bin/bash
########################################
# Presubmit script triggered by Prow.  #
########################################

WD=$(dirname $0)
WD=$(cd $WD; pwd)
ROOT=$(dirname $WD)
source spectro/prow/functions.sh

# Exit immediately for non zero status
set -e
# Check unset variables
set -u
# Print command trace
set -x



create_images

exit 0
