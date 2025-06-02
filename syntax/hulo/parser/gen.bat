@REM Copyright 2025 The Hulo Authors. All rights reserved.
@REM Use of this source code is governed by a MIT-style
@REM license that can be found in the LICENSE file.
@echo off
cd ./grammar
java -jar ..\antlr.jar -Dlanguage=Go -visitor -no-listener -o ..\generated -package generated *.g4
