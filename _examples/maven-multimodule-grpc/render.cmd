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

echo Generating Maven multi-module project...
"%RENDER%" "%SCRIPT_DIR%templates" "%SCRIPT_DIR%project.yaml" ^
    -o "%OUTPUT_DIR%"

echo.
echo Output generated in: %OUTPUT_DIR%
echo.
dir /s /b "%OUTPUT_DIR%" 2>nul
