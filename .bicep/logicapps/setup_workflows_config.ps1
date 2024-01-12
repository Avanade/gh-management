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

    echo "WORKFLOWS_APPSETTINGS=$workflowNamesDictionary" >> $GITHUB_ENV
}