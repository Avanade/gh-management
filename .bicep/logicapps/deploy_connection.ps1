<#
.SYNOPSIS
    Deploys Azure Logic App workflow connection.
.DESCRIPTION
    Deploys the workflow connection by adding the reference into the connections.json file
    that is stored in the associated fileshare.
.PARAMETER ResourceGroup
    The name of the resource group where the Storage account is located.
.PARAMETER StorageAccount
    The name of the Storage account where the File Share is located.
.PARAMETER Id
    The full resource ID of the connection.
.PARAMETER RuntimeUrl
    The full runtime URL of the connection.
.PARAMETER Api
    The managed API reference of the connection.
.INPUTS
    None.
.OUTPUTS
    None.
.EXAMPLE
    New-WorkflowConnection `
        -ResourceGroup "rg-orchestration-ts" `
        -StorageAccount "stsampleworkflowsts" `
        -Id "/subscriptions/12952a70-6abe-4cf3-880a-81ce65fdc63f/resourceGroups/rg-orchestration-ts/providers/Microsoft.Web/connections/con-storage-deadletter" `
        -RuntimeUrl "/subscriptions/12952a70-6abe-4cf3-880a-81ce65fdc63f/resourceGroups/rg-orchestration-ts/providers/Microsoft.Web/connections/con-storage-deadletter" `
        -Api "/subscriptions/12952a70-6abe-4cf3-880a-81ce65fdc63f/providers/Microsoft.Web/locations/westeurope/managedApis/azureblob"
#>
function New-WorkflowConnection {
    Param(
    [Parameter(Mandatory = $true)]
    $ResourceGroup,
    [Parameter(Mandatory = $true)]
    $StorageAccount,
    [Parameter(Mandatory = $true)]
    $Api,
    [Parameter(Mandatory = $true)]
    $Id,
    [Parameter(Mandatory = $true)]
    $RuntimeUrl
    )

    $ErrorActionPreference = "Stop"
    $WarningPreference = "Continue"

    $names =  $Id.Split('/')
    $name = $names[$names.length - 1]

    # Get current IP
    $ip = (Invoke-WebRequest -uri "http://ifconfig.me/ip").Content

    try {
        Write-Host "Deploying workflow connection '" -NoNewLine
        Write-Host $name -NoNewline -ForegroundColor Yellow
        Write-Host "'..." -NoNewline

        # Open firewall
        Add-AzStorageAccountNetworkRule -ResourceGroupName $ResourceGroup -Name $StorageAccount -IPAddressOrRange $ip | Out-Null

        # Connects the Azure context and sets the subscription.
        # New-RpicTenantConnection

        # Static values
        $directoryPath = "/site/wwwroot/"

        # Get the storage account context
        $ctx = (Get-AzStorageAccount -ResourceGroupName $ResourceGroup -Name $StorageAccount).Context

        # Get the file share
        $fsName = (Get-AZStorageShare -Context $ctx).Name

        # Download connection file
        $configPath = $directoryPath + "connections.json"
        try {
            Get-AzStorageFileContent -Context $ctx -ShareName $fsName -Path $configPath -Force
            Start-Sleep -Seconds 5
        } catch {
            # No such file, create it
            $newContent = @"
{
    "managedApiConnections": {
    }
}
"@
           Set-Content -Path "./connections.json" -Value $newContent
        }

        $config = Get-Content -Path "./connections.json" | ConvertFrom-Json
        $sectionName = ('$config.managedApiConnections."' + $name + '"')
        $section = Invoke-Expression $sectionName
        if ($null -eq $section) {
            # Section missing, add it
            $value = @"
    {
        "api": {
            "id": "$Api"
        },
        "authentication": {
            "type": "ManagedServiceIdentity"
        },
        "connection": {
            "id": "$Id"
        },
        "connectionRuntimeUrl": "$RuntimeUrl"
    }
"@
            $config.managedApiConnections | Add-Member -Name $name -Value (Convertfrom-Json $value) -MemberType NoteProperty

        } else {
            # Update section just in case
            $section.api.id = $Api
            $section.connection.id = $Id
            $section.connectionRuntimeUrl = $RuntimeUrl
        }

        # Save and upload file
        $config | ConvertTo-Json -Depth 100 | Out-File ./connections.json
        Set-AzStorageFileContent -Context $ctx -ShareName $fsName -Source ./connections.json  -Path $configPath  -Force
        Remove-Item ./connections.json -Force
    } finally {
       # Remove the firewall rule
       Remove-AzStorageAccountNetworkRule -ResourceGroupName $ResourceGroup -Name $StorageAccount -IPAddressOrRange $ip | Out-Null
    }
    Write-Host "Done!" -ForegroundColor Green
}