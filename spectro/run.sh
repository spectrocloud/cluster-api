#!/bin/bash

rm generated/*

kustomize build core/global > generated/core-global.yaml
kustomize build core/base > generated/core-base.yaml

kustomize build bootstrap/global > generated/bootstrap-global.yaml
kustomize build bootstrap/base > generated/bootstrap-base.yaml

kustomize build controlplane/global > generated/controlplane-global.yaml
kustomize build controlplane/base > generated/controlplane-base.yaml
