//+build !faker

package service

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	azureDevOpsDelivery "github.com/monitoror/monitoror/monitorable/azuredevops/delivery/http"
	azureDevOpsModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
	azureDevOpsRepository "github.com/monitoror/monitoror/monitorable/azuredevops/repository"
	azureDevOpsUsecase "github.com/monitoror/monitoror/monitorable/azuredevops/usecase"
	configDelivery "github.com/monitoror/monitoror/monitorable/config/delivery/http"
	configRepository "github.com/monitoror/monitoror/monitorable/config/repository"
	configUsecase "github.com/monitoror/monitoror/monitorable/config/usecase"
	"github.com/monitoror/monitoror/monitorable/github"
	githubDelivery "github.com/monitoror/monitoror/monitorable/github/delivery/http"
	githubModels "github.com/monitoror/monitoror/monitorable/github/models"
	githubRepository "github.com/monitoror/monitoror/monitorable/github/repository"
	githubUsecase "github.com/monitoror/monitoror/monitorable/github/usecase"
	"github.com/monitoror/monitoror/monitorable/http"
	httpDelivery "github.com/monitoror/monitoror/monitorable/http/delivery/http"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"
	httpRepository "github.com/monitoror/monitoror/monitorable/http/repository"
	httpUsecase "github.com/monitoror/monitoror/monitorable/http/usecase"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	_jenkinsDelivery "github.com/monitoror/monitoror/monitorable/jenkins/delivery/http"
	_jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	_jenkinsRepository "github.com/monitoror/monitoror/monitorable/jenkins/repository"
	_jenkinsUsecase "github.com/monitoror/monitoror/monitorable/jenkins/usecase"
	"github.com/monitoror/monitoror/monitorable/ping"
	pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	pingRepository "github.com/monitoror/monitoror/monitorable/ping/repository"
	pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	pingdomDelivery "github.com/monitoror/monitoror/monitorable/pingdom/delivery/http"
	pingdomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"
	pingdomRepository "github.com/monitoror/monitoror/monitorable/pingdom/repository"
	pingdomUsecase "github.com/monitoror/monitoror/monitorable/pingdom/usecase"
	"github.com/monitoror/monitoror/monitorable/port"
	portDelivery "github.com/monitoror/monitoror/monitorable/port/delivery/http"
	portModels "github.com/monitoror/monitoror/monitorable/port/models"
	portRepository "github.com/monitoror/monitoror/monitorable/port/repository"
	portUsecase "github.com/monitoror/monitoror/monitorable/port/usecase"
	"github.com/monitoror/monitoror/monitorable/travisci"
	travisciDelivery "github.com/monitoror/monitoror/monitorable/travisci/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorable/travisci/models"
	travisciRepository "github.com/monitoror/monitoror/monitorable/travisci/repository"
	travisciUsecase "github.com/monitoror/monitoror/monitorable/travisci/usecase"

	"github.com/monitoror/monitoror/monitorable/stripe"
	stripeDelivery "github.com/monitoror/monitoror/monitorable/stripe/delivery"
	stripeModels "github.com/monitoror/monitoror/monitorable/stripe/models"
	stripeRepository "github.com/monitoror/monitoror/monitorable/stripe/repository"
	stripeUsecase "github.com/monitoror/monitoror/monitorable/stripe/usecase"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/system"

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

func (s *Server) initApis() {
	// Api group definition
	s.api = s.Group("/api/v1")

	// ------------- INFO ------------- //
	s.registerInfo()

	// ------------- CONFIG ------------- //
	s.registerConfig()

	// ------------- TILES ------------- //
	for variant := range s.config.Monitorable.Ping {
		registerTile(s.registerPing, variant, system.IsRawSocketAvailable())
	}
	for variant := range s.config.Monitorable.Port {
		registerTile(s.registerPort, variant, true)
	}
	for variant := range s.config.Monitorable.HTTP {
		registerTile(s.registerHTTP, variant, true)
	}
	for variant, conf := range s.config.Monitorable.Pingdom {
		registerTile(s.registerPingdom, variant, conf.IsValid())
	}
	for variant, conf := range s.config.Monitorable.TravisCI {
		registerTile(s.registerTravisCI, variant, conf.IsValid())
	}
	for variant, conf := range s.config.Monitorable.Jenkins {
		registerTile(s.registerJenkins, variant, conf.IsValid())
	}
	for variant, conf := range s.config.Monitorable.AzureDevOps {
		registerTile(s.registerAzureDevOps, variant, conf.IsValid())
	}
	for variant, conf := range s.config.Monitorable.Github {
		registerTile(s.registerGithub, variant, conf.IsValid())
	}
	for variant, conf := range s.config.Monitorable.Stripe {
		registerTile(s.registerStripe, variant, conf.IsValid())
	}
}

func (s *Server) registerInfo() {
	infoDelivery := handlers.NewHTTPInfoDelivery()
	s.api.GET("/info", s.cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoDelivery.GetInfo))
}

func (s *Server) registerConfig() {
	repository := configRepository.NewConfigRepository()
	usecase := configUsecase.NewConfigUsecase(repository, s.store, s.config.DownstreamCacheExpiration)
	delivery := configDelivery.NewConfigDelivery(usecase)

	s.api.GET("/config", delivery.GetConfig)

	s.configHelper = usecase
}

func (s *Server) registerPing(variant string) {
	pingConfig := s.config.Monitorable.Ping[variant]

	repository := pingRepository.NewPingRepository(pingConfig)
	usecase := pingUsecase.NewPingUsecase(repository)
	delivery := pingDelivery.NewPingDelivery(usecase)

	// Register route to echo
	route := s.api.GET("/ping", s.cm.UpstreamCacheHandler(delivery.GetPing))

	// Register data for config hydration
	s.configHelper.RegisterTile(ping.PingTileType, &pingModels.PingParams{}, route.Path, pingConfig.InitialMaxDelay)
}

func (s *Server) registerPort(variant string) {
	portConfig := s.config.Monitorable.Port[variant]

	repository := portRepository.NewPortRepository(portConfig)
	usecase := portUsecase.NewPortUsecase(repository)
	delivery := portDelivery.NewPortDelivery(usecase)

	// Register route to echo
	route := s.api.GET("/port", s.cm.UpstreamCacheHandler(delivery.GetPort))

	// Register data for config hydration
	s.configHelper.RegisterTile(port.PortTileType, &portModels.PortParams{}, route.Path, portConfig.InitialMaxDelay)
}

func (s *Server) registerHTTP(variant string) {
	httpConfig := s.config.Monitorable.HTTP[variant]

	repository := httpRepository.NewHTTPRepository(httpConfig)
	usecase := httpUsecase.NewHTTPUsecase(repository, s.store, s.config.UpstreamCacheExpiration)
	delivery := httpDelivery.NewHTTPDelivery(usecase)

	// Register route to echo
	httpGroup := s.api.Group("/http")
	routeStatus := httpGroup.GET("/status", s.cm.UpstreamCacheHandler(delivery.GetHTTPStatus))
	routeRaw := httpGroup.GET("/raw", s.cm.UpstreamCacheHandler(delivery.GetHTTPRaw))
	routeJSON := httpGroup.GET("/formatted", s.cm.UpstreamCacheHandler(delivery.GetHTTPFormatted))

	// Register data for config hydration
	s.configHelper.RegisterTile(http.HTTPStatusTileType,
		&httpModels.HTTPStatusParams{}, routeStatus.Path, httpConfig.InitialMaxDelay)
	s.configHelper.RegisterTile(http.HTTPRawTileType,
		&httpModels.HTTPRawParams{}, routeRaw.Path, httpConfig.InitialMaxDelay)
	s.configHelper.RegisterTile(http.HTTPFormattedTileType,
		&httpModels.HTTPFormattedParams{}, routeJSON.Path, httpConfig.InitialMaxDelay)
}

func (s *Server) registerPingdom(variant string) {
	pingdomConfig := s.config.Monitorable.Pingdom[variant]

	repository := pingdomRepository.NewPingdomRepository(pingdomConfig)
	usecase := pingdomUsecase.NewPingdomUsecase(repository, pingdomConfig, s.store)
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// Register route to echo
	pingdomGroup := s.api.Group(fmt.Sprintf("/pingdom/%s", variant))
	route := pingdomGroup.GET("/check", s.cm.UpstreamCacheHandler(delivery.GetCheck))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(pingdom.PingdomCheckTileType,
		variant, &pingdomModels.CheckParams{}, route.Path, pingdomConfig.InitialMaxDelay)
	s.configHelper.RegisterDynamicTileWithConfigVariant(pingdom.PingdomChecksTileType,
		variant, &pingdomModels.ChecksParams{}, usecase)
}

func (s *Server) registerTravisCI(variant string) {
	travisCIConfig := s.config.Monitorable.TravisCI[variant]

	repository := travisciRepository.NewTravisCIRepository(travisCIConfig)
	usecase := travisciUsecase.NewTravisCIUsecase(repository)
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// Register route to echo
	travisCIGroup := s.api.Group(fmt.Sprintf("/travisci/%s", variant))
	route := travisCIGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(travisci.TravisCIBuildTileType,
		variant, &travisciModels.BuildParams{}, route.Path, travisCIConfig.InitialMaxDelay)
}

func (s *Server) registerJenkins(variant string) {
	jenkinsConfig := s.config.Monitorable.Jenkins[variant]

	repository := _jenkinsRepository.NewJenkinsRepository(jenkinsConfig)
	usecase := _jenkinsUsecase.NewJenkinsUsecase(repository)
	delivery := _jenkinsDelivery.NewJenkinsDelivery(usecase)

	// Register route to echo
	jenkinsGroup := s.api.Group(fmt.Sprintf("/jenkins/%s", variant))
	route := jenkinsGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(jenkins.JenkinsBuildTileType,
		variant, &_jenkinsModels.BuildParams{}, route.Path, jenkinsConfig.InitialMaxDelay)
	s.configHelper.RegisterDynamicTileWithConfigVariant(jenkins.JenkinsMultiBranchTileType,
		variant, &_jenkinsModels.MultiBranchParams{}, usecase)
}

func (s *Server) registerAzureDevOps(variant string) {
	azureDevOpsConfig := s.config.Monitorable.AzureDevOps[variant]

	repository := azureDevOpsRepository.NewAzureDevOpsRepository(azureDevOpsConfig)
	usecase := azureDevOpsUsecase.NewAzureDevOpsUsecase(repository)
	delivery := azureDevOpsDelivery.NewAzureDevOpsDelivery(usecase)

	// Register route to echo
	azureGroup := s.api.Group(fmt.Sprintf("/azuredevops/%s", variant))
	routeBuild := azureGroup.GET("/build", s.cm.UpstreamCacheHandler(delivery.GetBuild))
	routeRelease := azureGroup.GET("/release", s.cm.UpstreamCacheHandler(delivery.GetRelease))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(azuredevops.AzureDevOpsBuildTileType,
		variant, &azureDevOpsModels.BuildParams{}, routeBuild.Path, azureDevOpsConfig.InitialMaxDelay)
	s.configHelper.RegisterTileWithConfigVariant(azuredevops.AzureDevOpsReleaseTileType,
		variant, &azureDevOpsModels.ReleaseParams{}, routeRelease.Path, azureDevOpsConfig.InitialMaxDelay)
}

func (s *Server) registerGithub(variant string) {
	githubConfig := s.config.Monitorable.Github[variant]
	// Custom UpstreamCacheExpiration only for count because github has no-cache for this query and the rate limit is 30req/Hour
	countCacheExpiration := time.Millisecond * time.Duration(githubConfig.CountCacheExpiration)

	repository := githubRepository.NewGithubRepository(githubConfig)
	usecase := githubUsecase.NewGithubUsecase(repository)
	delivery := githubDelivery.NewGithubDelivery(usecase)

	// Register route to echo
	azureGroup := s.api.Group(fmt.Sprintf("/github/%s", variant))
	routeCount := azureGroup.GET("/count", s.cm.UpstreamCacheHandlerWithExpiration(countCacheExpiration, delivery.GetCount))
	routeChecks := azureGroup.GET("/checks", s.cm.UpstreamCacheHandler(delivery.GetChecks))

	// Register data for config hydration
	s.configHelper.RegisterTileWithConfigVariant(github.GithubCountTileType,
		variant, &githubModels.CountParams{}, routeCount.Path, githubConfig.InitialMaxDelay)
	s.configHelper.RegisterTileWithConfigVariant(github.GithubChecksTileType,
		variant, &githubModels.ChecksParams{}, routeChecks.Path, githubConfig.InitialMaxDelay)
	s.configHelper.RegisterDynamicTileWithConfigVariant(github.GithubPullRequestTileType,
		variant, &githubModels.PullRequestParams{}, usecase)
}

func (s *Server) registerStripe(variant string) {
	stripeConfig := s.config.Monitorable.Stripe[variant]
	countCacheExpiration := time.Millisecond * time.Duration(stripeConfig.CountCacheExpiration)

	repository := stripeRepository.NewStripeRepository(stripeConfig)
	usecase := stripeUsecase.NewStripeUsecase(repository)
	delivery := stripeDelivery.NewStripeDelivery(usecase)

	azureGroup := s.api.Group(fmt.Sprintf("/stripe/%s", variant))
	routeCount := azureGroup.GET("/count", s.cm.UpstreamCacheHandlerWithExpiration(countCacheExpiration, delivery.GetCount))

	s.configHelper.RegisterTileWithConfigVariant(stripe.StripeCountTileType,
		variant, &stripeModels.CountParams{}, routeCount.Path, stripeConfig.InitialMaxDelay)
}
