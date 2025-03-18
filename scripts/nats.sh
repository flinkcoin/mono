#!/bin/bash


docker pull nats:latest
docker run -p 4222:4222 -ti nats:latest


