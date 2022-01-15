#!/bin/bash

go test -race -covermode=atomic -coverprofile=coverage.out \
    github.com/SLedunois/b3lb/pkg/admin \
    github.com/SLedunois/b3lb/pkg/api \
    github.com/SLedunois/b3lb/pkg/app \
    github.com/SLedunois/b3lb/pkg/config \
    github.com/SLedunois/b3lb/pkg/utils \
    github.com/SLedunois/b3lb/pkg/restclient
