@echo off
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%~dp0install.ps1" %*
set "install_result=%errorlevel%"
if not "%install_result%"=="0" echo A instalacao encontrou um problema. Leia a mensagem acima.
pause
exit /b %install_result%
