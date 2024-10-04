@allowed([
  'test'
  'uat'
  'prod'
])
param activeEnv string

@secure()
param serverName string
@secure()
param administratorLogin string
@secure()
param administratorLoginPassword string

@secure()
param databaseName string

param skuName string = 'Basic'
param skuTier string = 'Basic'

param location string = resourceGroup().location

resource sqlServer 'Microsoft.Sql/servers@2022-08-01-preview' = {
  name: serverName
  location: location
  tags: {
    project: 'gh-management,Approval System'
    env: activeEnv == 'prod' ? 'prod' : 'test,uat'
  }
  properties: {
    administratorLogin: administratorLogin
    administratorLoginPassword: administratorLoginPassword
  }
}

resource sqlServerDatabase 'Microsoft.Sql/servers/databases@2022-08-01-preview' = {
  parent: sqlServer
  name: databaseName
  location: location
  tags: {
    project: 'gh-management'
    env: activeEnv
  }
  sku: {
    name: skuName
    tier: skuTier
  }
}
