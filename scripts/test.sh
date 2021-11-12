#!/bin/bash
if [ "$(docker images | grep sledunois/bbb-dev | wc -l)" -eq "0" ]; then
    sh -c "$(dirname "$0")/build_image.sh";
fi

go test -race -covermode=atomic -coverprofile=coverage.out b3lb/pkg/admin b3lb/pkg/api b3lb/pkg/app b3lb/pkg/config b3lb/pkg/utils
