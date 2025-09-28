#!/bin/bash

HOSTNAME=`hostname`

git add .
git commit -m "auth commit from ${HOSTNAME}"
git push github main
git push gitee main
