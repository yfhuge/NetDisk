package config

type System struct {
	Addr string `mapstructure:"addr"`
}

type MySQLConf struct {
	Host         string `mapstructure:"host"`
	UserName     string `mapstructure:"username"`
	Port         int    `mapstructure:"port"`
	DBName       string `mapstructure:"databaseName"`
	DBPwd        string `mapstructure:"dbPwd"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
	MaxOpenConns int    `mapstructre:"maxOpenConns"`
}

type RedisConf struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	PassWord string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type OssConf struct {
	Bucket          string `mapstructure:"bucket"`
	EndPoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	AccessKeySecret string `mapstructure:"accesskeySecret"`
}

type RabbitConf struct {
	AsyncTransferEnable  bool   `mapstructure:"asyncTransferEnable"`
	RabbitURL            string `mapstructure:"rabbitURL"`
	TransExchangeName    string `mapstructure:"transExchangeName"`
	TransOSSQueueName    string `mapstructure:"transOSSQueueName"`
	TransOSSErrQueueName string `mapstructure:"transOSSErrQueueName"`
	TransOSSRoutingKey   string `mapstructure:"transOSSRoutingKey"`
}

type Server struct {
	System   System     `mapstructure:"system"`
	DBConf   MySQLConf  `mapstructure:"dbConf"`
	RDBConf  RedisConf  `mapstructure:"rdbConf"`
	OSSConf  OssConf    `mapstructure:"ossConf"`
	RBMQConf RabbitConf `mapstructure:"rabbitConf"`
}
