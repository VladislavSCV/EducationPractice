@echo off
set local

set "go_file=./main.go"

REM Проверка установлени ли go
go version > nul 2>&1
if errorlevel 1 (
    echo Go не установлен( 
    exit /b
)

REM Запускаем файл
go run %go_file%

endlocal