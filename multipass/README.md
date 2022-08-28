# Canonical Multipass VM

Use this directory to setup a [Canonical Multipass](https://multipass.run/) Ubuntu VM to build and run project examples. This is useful if you want to build/run project in a non-Linux environment (i.e. Apple's OSX) or do not wish to test against your local environment directly.

## Pre-requisites

* Download [Canonical Multipass](https://multipass.run/)
* `envsubst` util command (i.e. `brew install gettext` on Mac)

## Run the VM

Use the `start.sh` script to launch the VM. Ensure that your local machine has the required spare CPUs and memory, otherwise, adjust accordingly. 

Once launched, use `multipass shell go4vldev` to log into the ubuntu VM.