@echo off
setlocal enabledelayedexpansion

set SCRIPT_DIR=%~dp0
set RENDER=%SCRIPT_DIR%..\..\render.exe
set OUTPUT_DIR=%SCRIPT_DIR%output

:: Build render if needed
if not exist "%RENDER%" (
    echo Building render...
    pushd %SCRIPT_DIR%..\..
    go build -o render.exe ./cmd/render
    popd
)

:: Clean and create output directory
if exist "%OUTPUT_DIR%" rmdir /s /q "%OUTPUT_DIR%"
mkdir "%OUTPUT_DIR%"

echo Generating Terraform modules...
echo.
echo Note: This example demonstrates templates for modules and environments.
echo Full generation requires combining module and environment data.
echo See README.md for details.
echo.

echo Module template preview (first 20 lines):
type "%SCRIPT_DIR%module-template\main.tf.tmpl" | findstr /n "^" | findstr "^[1-9]:" | findstr /v "^[2-9][0-9]:"

echo.
echo To generate, prepare per-module YAML files and run:
echo   render module-template\main.tf.tmpl vpc.yaml -o output\modules\vpc\main.tf
