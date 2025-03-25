@echo off
echo Starting AI-Powered Personal Knowledge Base...

echo Checking for C compiler...
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Warning: GCC compiler not found. Attempting to run without CGO...
    set CGO_ENABLED=0
) else (
    echo C compiler found, running with CGO enabled...
    set CGO_ENABLED=1
)

cd ..
go run cmd/server/main.go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Application exited with code %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
