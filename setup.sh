#!/bin/bash
# Copyright 2025 The Hulo Authors. All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.

echo Installing mage...
go install github.com/magefile/mage@latest
mage setup