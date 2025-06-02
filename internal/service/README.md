# internal/service
`service`はcoreやioなどのモジュール群をまとめて提供する概念。  
実体は主にcoreにあるが、そのロジックを利用するためのinterfaceのみを提供する。  
また、それらのモジュールを利用してservice interfaceを実装するオブジェクトを生成するのもここの役割。  
したがって、interfaceと、生成ロジックのみがここに実装されている。
