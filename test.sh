#! /usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


export FORGE_API_KEY=""
export FORGE_SERVER_ID=""
export FORGE_SITE_ID=""

# Run the test
go test -v ./internal/...
