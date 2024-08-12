#!/bin/bash

set -e
set -o pipefail

bin_name="dnsservices"

resolvers_file="resolvers.yaml"
signatures_file="signatures.yaml"

destination="/usr/local/bin/"

rm -rf $destination/$resolvers_file
rm -rf $destination/$signatures_file
rm -rf $destination/$bin_name

echo "Uninstalled successfully!"