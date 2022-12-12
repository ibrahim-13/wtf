@ECHO OFF
CLS

IF [%1] == [] GOTO InvalidArg

SET WTF_WIDGET_NAME=%1
go run generator\textwidget.go
EXIT

:InvalidArg
ECHO Invalid argument, pass widget name as the first arg