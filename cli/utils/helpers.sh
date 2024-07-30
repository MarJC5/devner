#!/bin/bash

wrap_text() {
    local text=$1
    local width=$2
    echo "$text" | fold -sw $width
}

wait_for_containers() {
    sleep 1
}