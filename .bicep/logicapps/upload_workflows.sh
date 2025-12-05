#!/bin/bash
set -e  # Exit on error

# Verify Azure login
echo "Verifying Azure login..."
az account show > /dev/null 2>&1 || { echo "ERROR: Not logged into Azure CLI"; exit 1; }
echo "Azure login verified."

# GET ACCESS KEY
echo "Retrieving storage account key..."
ACCOUNT_KEY=$(az storage account keys list -n "$1" --query '[0].value' -o tsv 2>&1)

if [ -z "$ACCOUNT_KEY" ] || [[ "$ACCOUNT_KEY" == *"ERROR"* ]]; then
    echo "ERROR: Failed to retrieve storage account key for $1"
    echo "Output: $ACCOUNT_KEY"
    exit 1
fi

echo "Storage account key retrieved successfully."

# UPLOAD ALL FILES INSIDE OF WORKFLOWS FOLDER
echo "Uploading workflows from $3 to $2..."
az storage file upload-batch --destination "$2" --source "$3" --account-name "$1" --account-key "$ACCOUNT_KEY" 2>&1 || {
    echo "ERROR: Upload failed. Trying with connection string instead..."
    # Try alternate method using connection string
    CONNECTION_STRING=$(az storage account show-connection-string -n "$1" --query connectionString -o tsv 2>&1)
    if [ $? -eq 0 ]; then
        az storage file upload-batch --destination "$2" --source "$3" --connection-string "$CONNECTION_STRING"
    else
        echo "ERROR: Both upload methods failed"
        exit 1
    fi
}

echo "Upload completed successfully."