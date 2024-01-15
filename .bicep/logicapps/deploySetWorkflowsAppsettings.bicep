param env string
param existingAppSettings object
param appSettings object

var resourceName = 'Ghmgm'
var logicAppName = '${resourceName}LA${toUpper(first(env))}${substring(env, 1)}'

resource LALogicApp 'Microsoft.Web/sites@2022-03-01' existing = {
  name: logicAppName
}

resource LALogicAppConfig 'Microsoft.Web/sites/config@2022-03-01' = {
  name: 'appsettings'
  parent: LALogicApp
  properties: union(existingAppSettings, appSettings)
}
