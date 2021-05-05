# cryco

[![Coverage Status](https://coveralls.io/repos/github/mengstr/cryco/badge.svg?branch=main)](https://coveralls.io/github/mengstr/cryco?branch=main)
[![codecov](https://codecov.io/gh/mengstr/cryco/branch/main/graph/badge.svg?token=NBE7FS5L7Z)](https://codecov.io/gh/mengstr/cryco)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/c6c13a3b3ce6421091113fce5fde24ef)](https://www.codacy.com/gh/mengstr/cryco/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mengstr/cryco&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/mengstr/cryco)](https://goreportcard.com/report/github.com/mengstr/cryco)
![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/mengstr/cryco?include_prereleases)

WIP - Golang package for handling encrypted config/settings files

## Values applied to struct

The fields in the struct will be receiving values from multiple sources. They are applied in the following order:

- Step 1) Values specified using the 'default' tag in the struct
- Step 2) Values from files
- Step 3) Values from environment variables

