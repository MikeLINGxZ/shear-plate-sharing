package main

func run() error {
	InitConfig()
	if config.Role == "server" {
		err := runServer()
		if err != nil {
			return err
		}
	}
	runClient()
	return nil
}
