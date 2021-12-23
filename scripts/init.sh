#!/bin/bash

go get -u github.com/gin-gonic/gin

# pre-commit installation
curl https://pre-commit.com/install-local.py | python -
source ~/.profile

pip install pre-commit

pre-commit --version
pre-commit install

npm install --save-dev @commitlint/{config-conventional,cli} husky
npx husky install
npx husky add .husky/commit-msg "npx --no -- commitlint --edit \"$1\""

sudo apt-get install tidy