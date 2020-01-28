#!/bin/bash

SERVICE="host-daemon"

go build -o $SERVICE

systemctl stop $SERVICE
systemctl disable $SERVICE.service

cp $SERVICE.service /lib/systemd/system/
cp host-daemon /bin/

systemctl enable $SERVICE.service

