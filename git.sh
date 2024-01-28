#!/bin/bash

#Check if the commit message was supplied
if [ -z "$1" ]
then
	echo "Please provide commit message"
	exit 1
fi

#Run git commands
git add .
git commit -m "$1"
git push