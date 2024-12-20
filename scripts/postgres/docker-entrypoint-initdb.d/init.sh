#! /usr/bin/env bash
set -e -o pipefail

psql -v ON_ERROR_STOP=1 --username root <<-EOSQL
	CREATE DATABASE exchange OWNER root;
EOSQL
