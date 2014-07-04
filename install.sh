#!/bin/bash
# Install script

export GOPATH=/root/go
echo "go install"
OUT="$(/usr/local/go/bin/go install 2>&1)"
#OUT=$(date)
if [[ -n "$OUT" ]]; then
	echo "failed"
	exit 1
fi

echo "ok"
echo "stop hidemyemail server"

OUT="$(kill $(ps -ef | awk ' /'bin\\/hidemyemail'/ {print $2}'))"

OUT="$(ps -ef | awk ' /'bin\\/hidemyemail'/ {print $2}')"
if [[ -n "$OUT" ]]; then
	echo "failed"	
fi
echo "ok"

echo "copy html"
yes | cp -frd html /root/go/bin/hidemyemail-res
echo "copy images"
yes | cp -frd images /root/go/bin/hidemyemail-res

nohup /root/go/bin/hidemyemail > /root/go/bin/hidemyemail.log &

