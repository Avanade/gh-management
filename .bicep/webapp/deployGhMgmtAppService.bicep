param appServicePlanName string

param location string = resourceGroup().location

param projectName string

param imageName string

@allowed([
  'test'
  'uat'
  'prod'
])
param activeEnv string

@secure()
param sqlServerName string

@secure()
param containerServer string

@secure()
param appServiceSettings object

@allowed([
  'F1'
  'B1'
  'P1v2'
  'P2v2'
  'P3v2'
  'P1V3'
  'P2V3'
  'P3V3'
])
param sku string = 'P1v2'

resource ghmgmtAppServicePlan 'Microsoft.Web/serverfarms@2020-06-01' = {
  name: appServicePlanName
  location: location
  properties: {
    reserved: true
  }
  sku: {
    name: sku
  }
  kind: 'linux'
}

var appServiceName = '${projectName}-${activeEnv}'

resource ghmgmtAppService 'Microsoft.Web/sites@2022-03-01' = {
  name: appServiceName
  location: location
  properties: {
    serverFarmId: ghmgmtAppServicePlan.id
    siteConfig: {
      appSettings: [for item in items(appServiceSettings): {
        name: item.key
        value: item.value
      }]
      linuxFxVersion: 'DOCKER|${containerServer}/${imageName}'
    }
  }
}

var possibleOutboundIpAddressesList = split(ghmgmtAppService.properties.possibleOutboundIpAddresses, ',')

module sqlServerFirewalls '../sql/sqlServerFirewallRules.bicep' = {
  name: 'ghmgmtSqlServerFirewalls'
  params: {
    outboundIpAddresses: possibleOutboundIpAddressesList
    projectName: appServiceName
    sqlServerName: sqlServerName
  }
}

module ghmgmtFrontDoor 'azureFrontDoor.bicep' = {
  name: 'frontdoor'
  params: {
    backendAddress: ghmgmtAppService.properties.defaultHostName
    frontDoorName: '${projectName}fd${activeEnv}'
  }
}

// TAGS
resource ghmgmtAppServicePlanTags 'Microsoft.Resources/tags@2022-09-01' = {
  name: 'default'
  scope: ghmgmtAppServicePlan
  properties: {
    tags: {
      project : 'gh-management,Approval System'
      env : 'test,uat,prod'
    }
  }
}

resource ghmgmtAppServiceTags 'Microsoft.Resources/tags@2022-09-01' = {
  name:  'default'
  scope: ghmgmtAppService
  properties: {
    tags: {
      project: 'gh-management'
      env: activeEnv
    }
  }
}
