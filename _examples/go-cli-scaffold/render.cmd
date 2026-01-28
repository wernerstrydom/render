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

echo === Pass 1: Generating base CLI structure ===
"%RENDER%" "%SCRIPT_DIR%templates" "%SCRIPT_DIR%cli.yaml" ^
    -o "%OUTPUT_DIR%"

echo.
echo === Pass 2: Generating command files ===
"%RENDER%" "%SCRIPT_DIR%command-template\command.go.tmpl" "%SCRIPT_DIR%commands.yaml" ^
    -o "%OUTPUT_DIR%\cmd\{{.name}}.go"

echo.
echo Output generated in: %OUTPUT_DIR%
echo.
dir /s /b "%OUTPUT_DIR%" 2>nul
