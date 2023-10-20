#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/alerts/types/expected_keepers.go -package mock -destination x/alerts/testutil/mocks.go