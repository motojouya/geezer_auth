package service

type ServiceLoader interface {
	JwtParserLoader
}

type serviceLoaderImpl struct {
	jwtParserLoaderImpl
}

func GetLoader() ServiceLoader {
	return &serviceLoaderImpl{
		jwtParserLoaderImpl: jwtParserLoaderImpl{},
	}
}
