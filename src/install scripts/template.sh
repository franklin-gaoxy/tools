#!/bin/bash

# System Info: Darwin
# Tested System: Darwin
# Service Info: scripts template
# Supported Command: install uninstall reinstall
# Desc: scripts template

# log print and write
function logs(){
    # check args
    if [ -z "$1" ] && [ -z "$2" ];then
        echo -e "\033[31m[ERROR]: The logs method requires two parameters!\033[0m"
        return 1
    fi

	if [[ "xinfo" == "x$1" ]];then
        echo -e "\033[34m[INFO]: $2\033[0m"
    elif [[ "xwarn" == "x$1" ]];then
        echo -e "\033[33m[WARN]: $2\033[0m"
    elif [[ "xerror" == "x$1" ]];then
        echo -e "\033[31m[ERROR]: $2\033[0m"
        exit 1
    elif [[ "xsuccess" == "x$1" ]];then
        echo -e "\033[32m[SUCCESS]: $2\033[0m"
    else
        echo -e "\033[33m[UNKNOWN]: Uncertain content!\033[0m"
    fi
}

# environment variable
function vars(){
    version="v0.0.0"
}


# Sub command check
function command(){
    if [ -z "$1" ];then
        logs "error" "Please enter the sub command:[install/uninstall/reinstall]"
    fi

    if [[ "x$1" == "xinstall" ]];then
        install
    elif [[ "x$1" == "xuninstall" ]];then
        uninstall
    elif [[ "x$1" == "xreinstll" ]];then
        uninstall
        logs "success" "uninstall success!start install ..."
        install
        logs "success" "installer success!"
    else
        logs "error" "unknown command!"
    fi
}

# Pre-check
function precheck(){
    # start check
    logs "info" "precheck success."
}

# Starting position or entrance of logical code
function main(){
    logs "info" "start run main func..."
}

function endcheck(){
    logs "info" "end check success."
}

# Start
logs
vars
command $1
precheck
main
endcheck