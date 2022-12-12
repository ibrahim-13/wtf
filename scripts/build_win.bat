@ECHO OFF
CLS

ECHO "____    __    ____ .___________. _______ "
ECHO "\   \  /  \  /   / |           ||   ____|"
ECHO " \   \/    \/   /  `---|  |----`|  |__   "
ECHO "  \            /       |  |     |   __|  "
ECHO "   \    /\    /        |  |     |  |     "
ECHO "    \__/  \__/         |__|     |__|     "

SET GOPROXY="https://proxy.golang.org,direct"
SET GO111MODULE=on

go build -v -o .\bin\wtfutil.exe