@ECHO OFF
IF "%1"=="NO-CACHE" docker build --no-cache --tag atlas-drg:latest .
IF NOT "%1"=="NO-CACHE" docker build --tag atlas-drg:latest .