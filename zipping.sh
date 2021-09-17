#!/bin/bash

if [ -e "operator.zip" ]; then
    rm operator.zip
fi


zip -r operator.zip . -X -x "*.DS_Store" -x "*.docker-secret" -x ".git"