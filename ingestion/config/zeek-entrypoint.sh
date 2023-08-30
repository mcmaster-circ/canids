#!/bin/sh

if [[ ! -z "${ZEEK_CMD}"]]; then
    echo "ZeekArgs = \"${ZEEK_CMD}\"" >> /usr/local/zeek/etc/zeekctl.cfg
fi

service cron start
zeekctl cron enable
zeekctl deploy
tail -f /dev/null