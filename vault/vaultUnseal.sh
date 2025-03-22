#!/bin/bash

echo "Starting vault unseal task"

# Load environment variables , unseal keys must be available in .env file
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Will unseal vault
docker exec -it vault vault operator unseal $Vault_unsealKey1
docker exec -it vault vault operator unseal $Vault_unsealKey2
docker exec -it vault vault operator unseal $Vault_unsealKey3