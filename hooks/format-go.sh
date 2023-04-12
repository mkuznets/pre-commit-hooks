#!/usr/bin/env sh

set -e # stop at first failure

local_prefix=""

while getopts ":l:" opt; do
  case ${opt} in
  l)
    local_prefix="$OPTARG"
    ;;
  \?)
    echo "Invalid option: -$OPTARG" 1>&2
    exit 1
    ;;
  esac
done

shift $((OPTIND - 1))

if [ $# -lt 1 ]; then
  echo "Usage: $0 <argument>"
  exit 1
fi

gofumpt -l -w "$@"

gci write "$@" \
  -s standard \
  ${local_prefix:+'-s'} ${local_prefix:+"prefix(${local_prefix})"} \
  -s default \
  -s blank \
  -s dot \
  --skip-generated
