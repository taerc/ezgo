#!/bin/sh
MESSAGE="AUTO-MSG"
SPACE=$(cd "$(dirname "$0")";pwd)
echo "${SPACE}"
cd ${SPACE}
git pull 
git add .
git commit -m ${MESSAGE}
git push

