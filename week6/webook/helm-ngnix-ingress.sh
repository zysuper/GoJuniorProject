#!/usr/bin/env bash

helm upgrade --install ingress-nginx ingress-nginx \
--repo https://kubernetes.github.io/ingress-nginx \
--namespace ingress-nginx --values override.yaml \
--create-namespace