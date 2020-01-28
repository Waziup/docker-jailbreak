# Docker Jailbreak

[![GoDoc](https://godoc.org/github.com/Waziup/docker-jailbreak?status.svg)](https://godoc.org/github.com/Waziup/docker-jailbreak)

Docker-Jailbreak is a daemon service the allows the execution of commands on the host machine from inside a docker container.

The host-daemon creates a unix socket `/var/run/host.sock` that can be mapped into a docker container to give the container access to host.

*Warning!* This service defeats the docker security!