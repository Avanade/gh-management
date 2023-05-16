param resourceName string = 'Ghmgm'
param env string
param location string = resourceGroup().location
param LAManageIdentityName string

var logicAppName = '${resourceName}LA${env}'
var fileShare = '${toLower(logicAppName)}fs${env}'
var accountKey = LAStorageAccount.listKeys().keys[0].value
var accountName = LAStorageAccount.name

resource LAManageIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' = {
  name: LAManageIdentityName
  location: location
}

resource LAStorageAccount 'Microsoft.Storage/storageAccounts@2022-09-01' = {
  name: toLower('${resourceName}sa')
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'Storage'
  properties: {
    minimumTlsVersion: 'TLS1_2'
    allowBlobPublicAccess: true
    networkAcls: {
      bypass: 'AzureServices'
      virtualNetworkRules: []
      ipRules: []
      defaultAction: 'Allow'
    }
    supportsHttpsTrafficOnly: true
    encryption: {
      services: {
        file: {
          keyType: 'Account'
          enabled: true
        }
        blob: {
          keyType: 'Account'
          enabled: true
        }
      }
      keySource: 'Microsoft.Storage'
    }
  }
}

resource LAAppServicePlan 'Microsoft.Web/serverfarms@2022-03-01' = {
  name: '${resourceName}ASP'
  location: location
  sku: {
    name: 'WS1'
    tier: 'WorkflowStandard'
    size: 'WS1'
    family: 'WS'
    capacity: 1
  }
  kind: 'elastic'
  properties: {
    perSiteScaling: false
    elasticScaleEnabled: true
    maximumElasticWorkerCount: 20
    isSpot: false
    reserved: false
    isXenon: false
    hyperV: false
    targetWorkerCount: 0
    targetWorkerSizeId: 0
    zoneRedundant: false
  }
}

resource LALogicApp 'Microsoft.Web/sites@2022-03-01' = {
  name: logicAppName
  location: location
  kind: 'functionapp,workflowapp'
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${LAManageIdentity.id}' : {}
    }
  }
  properties: {
    serverFarmId: LAAppServicePlan.id
    // clientAffinityEnabled: false
    // httpsOnly: true
  }
}

resource LALogicAppConfig 'Microsoft.Web/sites/config@2022-03-01' = {
  name: 'appsettings'
  parent: LALogicApp
  properties: {
    APP_KIND : 'workflowApp'
    AzureFunctionsJobHost__extensionBundle__id : 'Microsoft.Azure.Functions.ExtensionBundle.Workflows'
    AzureFunctionsJobHost__extensionBundle__version : '[1.*, 2.0.0)'
    AzureWebJobsStorage : 'DefaultEndpointsProtocol=https;AccountName=${accountName};AccountKey=${accountKey};EndpointSuffix=core.windows.net'
    FUNCTIONS_EXTENSION_VERSION : '~4'
    FUNCTIONS_WORKER_RUNTIME : 'node'
    WEBSITE_CONTENTAZUREFILECONNECTIONSTRING : 'DefaultEndpointsProtocol=https;AccountName=${accountName};AccountKey=${accountKey};EndpointSuffix=core.windows.net'
    WEBSITE_CONTENTSHARE : fileShare
    WEBSITE_NODE_DEFAULT_VERSION : '~14'
  }
}

output accountName string = LAStorageAccount.name
output destination string = '${fileShare}/site/wwwroot'
