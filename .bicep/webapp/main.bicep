param location string = resourceGroup().location

param projectName string

@secure()
param dockerImage string

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
  name: '${projectName}-asp'
  location: location
  properties: {
    reserved: true
  }
  sku: {
    name: sku
  }
  kind: 'linux'
}

resource ghmgmtAppService 'Microsoft.Web/sites@2022-03-01' = {
  name: projectName
  location: location
  properties: {
    serverFarmId: ghmgmtAppServicePlan.id
    siteConfig: {
      appSettings: [for item in items(appServiceSettings): {
        name: item.key
        value: item.value
      }]
      linuxFxVersion: 'DOCKER|${containerServer}/${dockerImage}'
    }
  }
}
