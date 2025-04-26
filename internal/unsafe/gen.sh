#!/bin/bash
# Copyright 2025 The Hulo Authors. All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.

java -jar ./antlr.jar -Dlanguage=Go -visitor -no-listener -package unsafe ./*.g4
