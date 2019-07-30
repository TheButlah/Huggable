#!/bin/bash
# Have this run every minute via crontab
chown -R huggable:huggable /opt/huggable.us
chmod -R ug+rw,o-rwx /opt/huggable.us
find /opt/huggable.us -type d -exec chmod ug+sx {} \;