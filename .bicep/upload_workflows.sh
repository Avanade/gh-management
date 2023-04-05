# GET ACCESS KEY
ACCOUNT_KEY=$(az storage account keys list -n $1 --query '[].value' -o tsv)

# UPLOAD ALL FILES INSIDE OF WORKFLOWS FOLDER
az storage file upload-batch --destination $2 --source ./.bicep/workflows --account-name $1 --account-key "$ACCOUNT_KEY"