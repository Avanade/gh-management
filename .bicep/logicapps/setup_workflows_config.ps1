function New-WokflowsConfig {
    Param(
    [Parameter(Mandatory = $true)]
    $WorkflowsPath
    )

    # Get the names of folders inside folders
    $workflowNames = Get-ChildItem -Path $WorkflowsPAth -Directory | Select-Object -ExpandProperty Name

    # Convert the names to an array
    $workflowNamesArray = $workflowNames -split ','

    $workflowNamesDictionary = @{}

    foreach ($workflowName in $workflowNamesArray) {
        $workflowNamesDictionary[$workflowName] = $false
    }

    $deploymentParameters = {
        "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentParameters.json#",
        "contentVersion": "1.0.0.0",
        "parameters": {
            "appSettings": {
                "value": $workflowNamesDictionary
            }
        }
    }

    $result = $deploymentParameters | ConvertTo-Json

    return $result
}