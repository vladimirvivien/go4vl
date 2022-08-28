#! /bin/bash

root_dir=$(dirname "${BASH_SOURCE[0]}")/..
SSH_PUBLIC_KEY=$(cat ~/.ssh/id_rsa.pub)
envsubst '$SSH_PUBLIC_KEY' < ./cloud-init.yaml | multipass launch jammy -v -n go4vldev --cpus 4 --mem 4g --disk 20g  --mount $root_dir:/home/ubuntu/go4vl  --cloud-init -
