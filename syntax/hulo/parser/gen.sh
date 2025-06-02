#!/bin/bash
# Copyright 2025 The Hulo Authors. All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.
cd ./grammar || exit 1
java -jar ../antlr.jar -Dlanguage=Go -visitor -no-listener -o ../generated -package generated *.g4
