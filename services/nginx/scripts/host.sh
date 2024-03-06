#!/bin/bash

set -e

echo "Please enter the host name:"
read host

# Add the domain to the /etc/hosts file if it doesn't already exist
if ! grep "$host" /etc/hosts; then
    ## Return line
    sudo -- sh -c -e "printf '\n' >> /etc/hosts"
    sudo -- sh -c -e "printf '127.0.0.1\t%-30s\n' ${host} >> /etc/hosts" # IPv4
    sudo -- sh -c -e "printf '::1\t\t%-30s\n' ${host} >> /etc/hosts" # IPv6
    sudo -- sh -c -e "printf '127.0.0.1\t%-30s\n' www.${host} >> /etc/hosts" # IPv4
    sudo -- sh -c -e "printf '::1\t\t%-30s\n' www.${host} >> /etc/hosts" # IPv6
    cat /etc/hosts
else
    echo "$host Host already existing in /etc/hosts"
fi