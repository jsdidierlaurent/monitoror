export enum TileType {
  HttpStatus = 'HTTP-STATUS',
  HttpRaw = 'HTTP-RAW',
  HttpFormatted = 'HTTP-FORMATTED',
  Ping = 'PING',
  Port = 'PORT',
  Pingdom = 'PINGDOM-CHECK',
  GitHubChecks = 'GITHUB-CHECKS',
  GitHubCount = 'GITHUB-COUNT',
  GitLab = 'GITLAB-BUILD',
  Travis = 'TRAVISCI-BUILD',
  Jenkins = 'JENKINS-BUILD',
  AzureDevOpsBuild = 'AZUREDEVOPS-BUILD',
  AzureDevOpsRelease = 'AZUREDEVOPS-RELEASE',
  StripeCount = 'STRIPE-COUNT',

  Empty = 'EMPTY',
  Group = 'GROUP',
}

export default TileType
