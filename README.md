# cryco

[![Coverage Status](https://coveralls.io/repos/github/mengstr/cryco/badge.svg?branch=main)](https://coveralls.io/github/mengstr/cryco?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mengstr/cryco)](https://goreportcard.com/report/github.com/mengstr/cryco)
![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/mengstr/cryco?include_prereleases)

WIP - Golang package for handling encrypted config/settings files

## Values applied to struct

The fields in the struct will be receiving values from multiple sources. They are applied in the following order:

- Step 1) Values specified using the 'default' tag in the struct
- Step 2) Values from files
- Step 3) Values from environment variables

