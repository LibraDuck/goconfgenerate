package config
import (
	"encoding/json"
	"io/ioutil"
)

type RedisConfig struct {
	Host	string	`json:"host"`
	Port	string	`json:"port"`
	Pass	string	`json:"pass"`
	Db	int	`json:"db"`
	
}

type MysqlConfig struct {
	Host	string	`json:"host"`
	Port	string	`json:"port"`
	User	string	`json:"user"`
	Pass	string	`json:"pass"`
	Dbname	string	`json:"dbname"`
	
}

type NatsConfig struct {
	Host	string	`json:"host"`
	Port	string	`json:"port"`
	User	string	`json:"user"`
	Pass	string	`json:"pass"`
	Topic	*TopicConfig	`json:"topic"`
	
}

type RobotConfig struct {
	Redis_Heart	string	`json:"redis_heart"`
	Heart_Interval	int	`json:"heart_interval"`
	
}

type DaemonConfig struct {
	Login_Overtime	int	`json:"login_overtime"`
	Scan_Overtime	int	`json:"scan_overtime"`
	Heart_Overtime	int	`json:"heart_overtime"`
	
}

type TopicConfig struct {
	Login	string	`json:"login"`
	Scan	string	`json:"scan"`
	Heart	string	`json:"heart"`
	Exit	string	`json:"exit"`
	
}

type GlobalConfig struct {
	Http_Port	int	`json:"http_port"`
	Redis	*RedisConfig	`json:"redis"`
	Mysql	*MysqlConfig	`json:"mysql"`
	Nats	*NatsConfig	`json:"nats"`
	Robot	*RobotConfig	`json:"robot"`
	Daemon	*DaemonConfig	`json:"daemon"`
	Docker	*DockerConfig	`json:"docker"`
	
}

var (
	config *GlobalConfig
)

func Config() *GlobalConfig {
	return config
}

//json格式的配置文件解析
func InitConfig(dir string) error {
	buf, err := ioutil.ReadFile(dir)
	if err != nil {
		return err
	}
	var c GlobalConfig
	if err := json.Unmarshal(buf, &c); err != nil {
		return err
	}
	config = &c
	return nil
}
