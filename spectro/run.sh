#!/bin/bash


#rm generated/*

kustomize build --load_restrictor none core/global > generated/core-global.yaml
kustomize build --load_restrictor none core/base > generated/core-base.yaml

kustomize build --load_restrictor none bootstrap/global > generated/bootstrap-global.yaml
kustomize build --load_restrictor none bootstrap/base > generated/bootstrap-base.yaml

kustomize build --load_restrictor none controlplane/global > generated/controlplane-global.yaml
kustomize build --load_restrictor none controlplane/base > generated/controlplane-base.yaml
