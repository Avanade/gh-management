param frontDoorName string
param backendAddress string
param customDomain string = ''

param activeEnv string

var withCustomDomain = !empty(customDomain)

var defaultFrontEndEndpointName = 'azurefd-net'
var customFrontEndEndpointName = 'custom-domain'

var frontEndEndpointName = withCustomDomain ? customFrontEndEndpointName : defaultFrontEndEndpointName

var frontEndEndpoints = [{
  name: defaultFrontEndEndpointName
  properties: {
    hostName: '${frontDoorName}.azurefd.net'
    sessionAffinityEnabledState: 'Disabled'
  }
}
{
  name: customFrontEndEndpointName
  properties: {
    hostName: customDomain
    sessionAffinityEnabledState: 'Disabled'
  }
}]

var loadBalancingSettingsName = 'loadBalancingSettings'
var healthProbeSettingsName = 'healthProbeSettings'
var routingRuleName = 'routingRule'
var backendPoolName = 'backendPool'

resource frontDoor 'Microsoft.Network/frontDoors@2021-06-01' = {
  name: frontDoorName
  location: 'global'
  properties: {
    enabledState: 'Enabled'
    frontendEndpoints: withCustomDomain ? frontEndEndpoints : [frontEndEndpoints[0]]
    loadBalancingSettings: [
      {
        name: loadBalancingSettingsName
        properties: {
          sampleSize: 4
          successfulSamplesRequired: 2
        }
      }
    ]

    healthProbeSettings: [
      {
        name: healthProbeSettingsName
        properties: {
          path: '/'
          protocol: 'Https'
          intervalInSeconds: 30
          enabledState: 'Enabled'
          healthProbeMethod: 'Head'
        }
      }
    ]
    backendPools: [
      {
        name: backendPoolName
        properties: {
          backends: [
            {
              address: backendAddress
              backendHostHeader: backendAddress
              httpPort: 80
              httpsPort: 443
              priority: 1
              weight: 50
              enabledState: 'Enabled'
            }
          ]
          loadBalancingSettings: {
            id: resourceId('Microsoft.Network/frontDoors/loadBalancingSettings', frontDoorName, loadBalancingSettingsName)
          }
          healthProbeSettings: {
            id: resourceId('Microsoft.Network/frontDoors/healthProbeSettings', frontDoorName, healthProbeSettingsName)
          }
        }
      }
    ]
    routingRules: [
      {
        name: routingRuleName
        properties: {
          frontendEndpoints: [
            {
              id: resourceId('Microsoft.Network/frontDoors/frontEndEndpoints', frontDoorName, frontEndEndpointName)
            }
          ]
          acceptedProtocols: [
            'Http'
            'Https'
          ]
          patternsToMatch: [
            '/*'
          ]
          routeConfiguration: {
            '@odata.type': '#Microsoft.Azure.FrontDoor.Models.FrontdoorForwardingConfiguration'
            forwardingProtocol: 'MatchRequest'
            backendPool: {
              id: resourceId('Microsoft.Network/frontDoors/backEndPools', frontDoorName, backendPoolName)
            }
          }
          enabledState: 'Enabled'
        }
      }
    ]
  }
}

//TAGS
resource frontDoorTags 'Microsoft.Resources/tags@2022-09-01' = {
  name: 'default'
  scope: frontDoor
  properties: {
    tags: {
      env : activeEnv
    }
  }
}
