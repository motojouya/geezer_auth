package service

type ServiceLoader interface {
	AuthorizationLoader
	JwtHandlerLoader
}

type serviceLoaderImpl struct {
	authorizationLoaderImpl
	jwtHandlerLoaderImpl
}

func GetLoader() ServiceLoader {
	return &serviceLoaderImpl{
		authorizationLoaderImpl: authorizationLoaderImpl{},
		jwtHandlerLoaderImpl:    jwtHandlerLoaderImpl{},
	}
}
