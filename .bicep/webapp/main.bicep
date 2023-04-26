param location string = resourceGroup().location

@secure()
param dockerImage string

@secure()
param containerServer string

@secure()
param containerUsername string

@secure()
param containerPassword string

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
  name: '${dockerImage}-asp'
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
  name: dockerImage
  location: location
  properties: {
    serverFarmId: ghmgmtAppServicePlan.id
    siteConfig: {
      appSettings: [
        {
          name: 'WEBSITES_ENABLE_APP_SERVICE_STORAGE'
          value: 'false'
        }
        {
          name: 'DOCKER_REGISTRY_SERVER_URL'
          value: 'https://${containerServer}'
        }
        {
          name: 'DOCKER_REGISTRY_SERVER_USERNAME'
          value: containerUsername
        }
        {
          name: 'DOCKER_REGISTRY_SERVER_PASSWORD'
          value: containerPassword
        }
      ]
      linuxFxVersion: 'DOCKER|${containerServer}/${dockerImage}'
    }
  }
}
