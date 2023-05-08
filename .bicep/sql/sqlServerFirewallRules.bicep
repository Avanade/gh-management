param projectName string
param sqlServerName string
param outboundIpAddresses array

resource ghmgmtSqlServer 'Microsoft.Sql/servers@2022-08-01-preview' existing = {
  name: sqlServerName
}

resource ghmgmtSqlServerFirewalls 'Microsoft.Sql/servers/firewallRules@2022-08-01-preview' = [for outboundIpAddress in outboundIpAddresses: {
  name: '${projectName}-${outboundIpAddress}'
  parent: ghmgmtSqlServer
  properties: {
    endIpAddress: outboundIpAddress
    startIpAddress: outboundIpAddress
  }
}]
