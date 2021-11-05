#!/bin/bash

sudo apt install httpie

go get -u github.com/gin-gonic/gin

# pre-commit installation
curl https://pre-commit.com/install-local.py | python -
source ~/.profile

pre-commit --version
pre-commit install
