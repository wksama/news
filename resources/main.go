package resources

func Init() {
	configInit()
	logInit()
	databaseInit()
	redisInit()
}
