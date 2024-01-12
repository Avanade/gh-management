param env string
param appSettings object

var resourceName = 'Ghmgm'
var logicAppName = '${resourceName}LA${toUpper(first(env))}${substring(env, 1)}'

resource LALogicApp 'Microsoft.Web/sites@2022-03-01' existing = {
  name: logicAppName
}

var currentAppSettings = list(resourceId('Microsoft.Web/sites/config', LALogicApp.name, 'appsettings'), '2022-03-01').properties


resource LALogicAppConfig 'Microsoft.Web/sites/config@2022-03-01' = {
  name: 'appsettings'
  parent: LALogicApp
  properties: union(currentAppSettings, appSettings)
}
