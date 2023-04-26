param location string = resourceGroup().location
param containerRegName string

resource containerRegistry 'Microsoft.ContainerRegistry/registries@2022-02-01-preview' existing = {
  name: containerRegName
}

var prefix = 'bps${uniqueString(resourceGroup().id)}'

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

resource bpsAppServicePlan 'Microsoft.Web/serverfarms@2020-06-01' = {
  name: '${prefix}AppServicePlan'
  location: location
  properties: {
    reserved: true
  }
  sku: {
    name: sku
  }
  kind: 'linux'
}

@description('The relative docker image name.')
param dockerImage string

resource bpsAppService 'Microsoft.Web/sites@2022-03-01' = {
  name: '${prefix}AppService'
  location: location
  properties: {
    serverFarmId: bpsAppServicePlan.id
    siteConfig: {
      appSettings: [
        {
          name: 'WEBSITES_ENABLE_APP_SERVICE_STORAGE'
          value: 'false'
        }
        {
          name: 'DOCKER_REGISTRY_SERVER_URL'
          value: 'https://${containerRegistry.properties.loginServer}'
        }
        {
          name: 'DOCKER_REGISTRY_SERVER_USERNAME'
          value: containerRegistry.name      }
        {
          name: 'DOCKER_REGISTRY_SERVER_PASSWORD'
          value: containerRegistry.listCredentials().passwords[0].value
        }
      ]
      linuxFxVersion: 'DOCKER|${containerRegistry.properties.loginServer}/${dockerImage}'
    }
  }
}

resource publishingcreds 'Microsoft.Web/sites/config@2021-01-01' existing = {
  parent: bpsAppService
  name: 'publishingcredentials'
}

var creds = list(publishingcreds.id, publishingcreds.apiVersion).properties.scmUri

resource containerRegistryWebhook 'Microsoft.ContainerRegistry/registries/webhooks@2022-02-01-preview' = {
  name: 'acrwebhook'
  location: location
  parent: containerRegistry
  properties: {
    actions: ['push']
    serviceUri: '${creds}/api/registry/webhook'
  }
}
