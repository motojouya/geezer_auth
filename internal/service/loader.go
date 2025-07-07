package service

type ServiceLoader interface {
	AuthorizationLoader
	JwtHandlerLoader
	DatabaseLoader
}

type serviceLoaderImpl struct {
	authorizationLoaderImpl
	jwtHandlerLoaderImpl
	databaseLoaderImpl
}

func GetLoader() ServiceLoader {
	return &serviceLoaderImpl{
		authorizationLoaderImpl: authorizationLoaderImpl{},
		jwtHandlerLoaderImpl:    jwtHandlerLoaderImpl{},
		databaseLoaderImpl:      databaseLoaderImpl{},
	}
}
