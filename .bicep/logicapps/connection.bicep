param resourceName string = 'Ghmgm'
param env string

param storageAccountName string
param logicAppName string
param location string = resourceGroup().location
param principalId string

// Get parent storage account
resource storage_account 'Microsoft.Storage/storageAccounts@2021-06-01' existing = {
  name: storageAccountName
}

//  Create connection
param connectionName string = '${resourceName}CN${toUpper(first(env))}${substring(env, 1)}'
resource connection 'Microsoft.Web/connections@2016-06-01' = {
  name: connectionName
  location: location
  kind: 'V2'
  properties: {
    displayName: connectionName
    api: {
      displayName: 'Azure Queues connection"'
      description: 'Azure Queue storage provides cloud messaging between application components. Queue storage also supports managing asynchronous tasks and building process work flows.'
      id:subscriptionResourceId('Microsoft.Web/locations/managedApis', location, 'azurequeues')
      type: 'Microsoft.Web/locations/managedApis'
    }
    parameterValues: {
      storageaccount: storage_account.name
      sharedkey: storage_account.listKeys().keys[0].value
    }
  }
}

// Create access policy for the connection
// Type not in Bicep yet but works fine
resource ConnectionPolicy 'Microsoft.Web/connections/accessPolicies@2016-06-01' = {
  name: '${connectionName}/${logicAppName}'
  location: location
  properties: {
    principal: {
      type: 'ActiveDirectory'
      identity: {
        tenantId: subscription().tenantId
        objectId: principalId
      }
    }
  }
  dependsOn: [
    connection
  ]
}

// Return the connection runtime URL, this needs to be set in the connection JSON file later
output connectionRuntimeUrl string = reference(connection.id, connection.apiVersion, 'full').properties.connectionRuntimeUrl
output api string = subscriptionResourceId('Microsoft.Web/locations/managedApis', location, 'azurequeues')
output id string = connection.id
output name string = connection.name
