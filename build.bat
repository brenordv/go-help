@echo off
echo Building Go applications for Linux and Windows...

:: Check if Go is installed
go version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Go is not installed or not in the PATH.
    exit /b 1
)

:: Create output directories
mkdir dist\linux >nul 2>&1
mkdir dist\windows >nul 2>&1

:: Set GOOS and GOARCH for Linux builds
set GOOS=linux
set GOARCH=amd64

:: Build cat.go into cat for Linux
echo Building cat for Linux...
go build -o dist\linux\cat cat.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build cat for Linux
    exit /b 1
)

:: Build touch.go into touch for Linux
echo Building touch for Linux...
go build -o dist\linux\touch touch.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build touch for Linux
    exit /b 1
)

:: Build guid.go into guid for Linux
echo Building guid for Linux...
go build -o dist\linux\guid guid.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build guid for Linux
    exit /b 1
)

:: Build ts.go into ts for Linux
echo Building ts for Linux...
go build -o dist\linux\ts ts.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build ts for Linux
    exit /b 1
)

:: Build ncsv.go into ncsv for Linux
echo Building ncsv for Linux...
go build -o dist\linux\ncsv ncsv.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build ncsv for Linux
    exit /b 1
)

:: Set GOOS and GOARCH for Windows builds
set GOOS=windows
set GOARCH=amd64

:: Build cat.go into cat.exe for Windows
echo Building cat for Windows...
go build -o dist\windows\cat.exe cat.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build cat for Windows
    exit /b 1
)

:: Build touch.go into touch.exe for Windows
echo Building touch for Windows...
go build -o dist\windows\touch.exe touch.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build touch for Windows
    exit /b 1
)

:: Build guid.go into guid.exe for Windows
echo Building guid for Windows...
go build -o dist\windows\guid.exe guid.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build guid for Windows
    exit /b 1
)

:: Build ts.go into ts.exe for Windows
echo Building ts for Windows...
go build -o dist\windows\ts.exe ts.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build ts for Windows
    exit /b 1
)

:: Build ncsv.go into ncsv.exe for Windows
echo Building ncsv for Windows...
go build -o dist\windows\ncsv.exe ncsv.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build ncsv for Windows
    exit /b 1
)

echo Build completed successfully for Linux and Windows.
exit /b 0
