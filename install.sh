#!/bin/bash

set -e
set -o pipefail

bin_name="dnscovery"

resolvers_file="resolvers.yaml"
signatures_file="signatures.yaml"

build_dir=".build/binaries/"
destination="/usr/local/bin/"

get_os_arch() {
  # Get the OS name
  os_name=$(uname -s | tr '[:upper:]' '[:lower:]')

  # Get the architecture
  arch_name=$(uname -m)

  # Convert architecture name to the desired format
  case "$arch_name" in
    x86_64)
      arch_name="amd64";;
    aarch64)
      arch_name="arm64";;
    armv7l)
      arch_name="arm";;
    i686)
      arch_name="386";;
    *)
      arch_name="windows";;
  esac

  echo "${os_name}-${arch_name}"
}

result=$(get_os_arch)

cp $resolvers_file $destination
cp $signatures_file $destination

cd $build_dir
file=$(ls | grep "$result")
echo "Detected arch: $file"
cp $file $bin_name

echo "Moving to $destination"
mv $bin_name $destination

echo "Installed successfully!"