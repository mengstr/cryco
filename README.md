# cryco

[![Coverage Status](https://coveralls.io/repos/github/mengstr/cryco/badge.svg?branch=main)](https://coveralls.io/github/mengstr/cryco?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mengstr/cryco)](https://goreportcard.com/report/github.com/mengstr/cryco)
[![Release](https://img.shields.io/github/release/mengstr/cryco.svg?label=Release)](https://github.com/mengstr/cryco/releases)

WIP - Golang package for handling encrypted config/settings files

## Values applied to struct

The fields in the struct will be receiving values from multiple sources. They are applied in the following order:

- Step 1) Values specified using the 'default' tag in the struct
- Step 2) Values from files
- Step 3) Values from environment variables

