param appServicePlanName string

param location string = resourceGroup().location

param projectName string

param imageName string

param runDeployFrontDoor bool
param frontDoorCustomDomain string

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

var appSettings = [for item in items(appServiceSettings): {
  name: item.key
  value: item.value
}]

resource ghmgmtAppService 'Microsoft.Web/sites@2022-03-01' = {
  name: appServiceName
  location: location
  properties: {
    serverFarmId: ghmgmtAppServicePlan.id
    siteConfig: {
      appSettings: union([{
        name: 'APPINSIGHTS_INSTRUMENTATIONKEY'
        value: appInsights.properties.InstrumentationKey
      }], appSettings)
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

var hostNameSslStates = filter(
  ghmgmtAppService.properties.hostNameSslStates, e => e.sslState == 'SniEnabled'
)

var backendAddress = length(hostNameSslStates) > 0 ? first (hostNameSslStates)!.name : ghmgmtAppService.properties.defaultHostName

module ghmgmtFrontDoor 'deployFrontDoor.bicep' = if(runDeployFrontDoor) {
  name: 'frontdoor'
  params: {
    backendAddress: backendAddress
    frontDoorName: '${projectName}fd-${activeEnv}'
    customDomain: frontDoorCustomDomain
    activeEnv: activeEnv
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

var appInsightName = toLower('${projectName}-${activeEnv}-appinsights')
var logAnalyticsName = toLower('${projectName}-${activeEnv}-loganalytics')

resource appInsights 'Microsoft.Insights/components@2020-02-02' = {
  name: appInsightName
  location: location
  kind: 'string'
  tags: {
    displayName: 'AppInsight'
    ProjectName: projectName
  }
  properties: {
    Application_Type: 'web'
    WorkspaceResourceId: logAnalyticsWorkspace.id
  }
}

resource logAnalyticsWorkspace 'Microsoft.OperationalInsights/workspaces@2020-08-01' = {
  name: logAnalyticsName
  location: location
  tags: {
    displayName: 'Log Analytics'
    ProjectName: projectName
  }
  properties: {
    sku: {
      name: 'PerGB2018'
    }
    retentionInDays: 120
    features: {
      searchVersion: 1
      legacy: 0
      enableLogAccessUsingOnlyResourcePermissions: true
    }
  }
}
