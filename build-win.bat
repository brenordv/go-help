@echo off
echo Building Go applications...

:: Check if Go is installed
go version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Go is not installed or not in the PATH.
    exit /b 1
)

:: Build cat.go into cat.exe
echo Building cat.exe...
go build -o cat.exe cat.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build cat.exe
    exit /b 1
)

:: Build touch.go into touch.exe
echo Building touch.exe...
go build -o touch.exe touch.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build touch.exe
    exit /b 1
)

:: Build guid.go into guid.exe
echo Building guid.exe...
go build -o guid.exe guid.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build guid.exe
    exit /b 1
)

:: Build ts.go into ts.exe
echo Building ts.exe...
go build -o ts.exe ts.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build ts.exe
    exit /b 1
)

echo Build completed successfully.
exit /b 0
