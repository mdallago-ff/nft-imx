#!/bin/bash
xunit-viewer -r $1 -o $2

if command -v open &> /dev/null
then
    open $2
    exit
fi

if command -v xdg-open &> /dev/null
then
    xdg-open $2
    exit
fi
