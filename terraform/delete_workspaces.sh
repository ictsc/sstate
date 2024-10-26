#!/bin/bash

terraform workspace select default || terraform workspace new default

for workspace in $(terraform workspace list | sed 's/^[* ]*//'); do
  if [ "$workspace" != "default" ]; then
    echo "Deleting workspace: $workspace"
    terraform workspace select "$workspace"
    terraform workspace select default
    terraform workspace delete "$workspace"
  fi
done
