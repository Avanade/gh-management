# GET ACCESS KEY
ACCOUNT_KEY=$(az storage account keys list -n $1 --query '[].value' -o tsv)

# UPLOAD ALL FILES INSIDE OF WORKFLOWS FOLDER
az storage file upload-batch --destination $2 --source $3 --account-name $1 --account-key "$ACCOUNT_KEY"

# Get the names of folders inside folders
$workflowNames = Get-ChildItem -Path $3 -Directory | Select-Object -ExpandProperty Name

# Convert the names to an array
$workflowNamesArray = $workflowNames -split ','

$workflowNamesDictionary = @{}

foreach ($workflowName in $workflowNamesArray) {
    $workflowNamesDictionary[$workflowName] = $false
}

echo "WORKFLOWS_APPSETTINGS=$workflowNamesDictionary" >> $GITHUB_ENV