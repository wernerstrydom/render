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

echo Generating Makefile for Go workspace...
"%RENDER%" "%SCRIPT_DIR%templates\Makefile.tmpl" "%SCRIPT_DIR%workspace.yaml" ^
    -o "%OUTPUT_DIR%\Makefile"

echo.
echo Output generated in: %OUTPUT_DIR%
echo.
dir "%OUTPUT_DIR%"
