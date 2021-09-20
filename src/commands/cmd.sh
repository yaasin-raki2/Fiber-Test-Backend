#!/bin/bash

RED="\033[1;31m"
PURPLE="\033[1;35m"
CYAN="\033[1;36m"
WHITE="\033[1;37m"
RESET="\033[0m"
ITALY="\e[3m"
UNDERSCORE="\e[4m"

if [ $# -eq 0 ] ; then
    echo -e "${RED}Please enter one or more commands to run from src${RESET}"
    echo -e "${WHITE}You can check available commands and access them using the --help or -h flags"
    exit 1
fi

CMDS=("populateUsers")
while [ -n "$1" ] ; do 
    if [ "$1" == "--help" ] || [ "$1" == "-h" ] ; then
        echo -e "${RED}Available Commands:${RESET}"
        C=0
        while [ "$C" -lt "${#CMDS[@]}" ] ; do
            if [ ! "$C" -eq  $(expr ${#CMDS[@]} - 1) ] ; then
                echo -e "${UNDERSCORE}${ITALY}${WHITE}\t${CMDS[$C]%?}"
            else
                echo -e "${UNDERSCORE}${ITALY}${WHITE}\t${CMDS[$C]}${RESET}" 
            fi
            (( C++ ))
        done
    elif [[ " ${CMDS[*]} " =~ " $1 " ]]; then
        echo -e "${PURPLE}Running $1 ...${RESET}"
        $(docker compose exec backend /bin/bash -c "go run src/commands/src/$1.go")
        echo -e "${CYAN}Command $1 executed successfully${RESET}"
    else
        echo -e "${RED}fuck u${RESET}"
        exit 1
    fi
    shift
done

exit 0