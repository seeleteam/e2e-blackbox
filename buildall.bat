@echo off
goto comment
    Build the command lines and tests in Windows.
    Must install gcc tool before building.
:comment

set para=%*
if not defined para (
    set para=all
)

for %%i in (%para%) do (
    call :%%i
)
pause
goto:eof

:all
call :run
goto:eof


:run
echo on
go build -o ./build/run.exe ./run/
@echo "Done run building"
@echo off
goto:eof


:clean
del build\* /q /f /s
@echo "Done clean the build dir"
@echo off
goto:eof
