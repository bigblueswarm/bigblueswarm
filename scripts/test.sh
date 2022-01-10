#!/bin/bash

go test -race -covermode=atomic -coverprofile=coverage.out b3lb/pkg/admin b3lb/pkg/api b3lb/pkg/app b3lb/pkg/config b3lb/pkg/utils b3lb/pkg/restclient
