package ssdb

type Config struct {

	// Database server hostname or IP. Leave blank if using unix sockets
	Host string `json:"host"`

	// Database server port. Leave blank if using unix sockets
	Port int `json:"port"`

	// TODO Password for authentication
	// Pass string `json:"pass"`

	// TODO A path of a UNIX socket file. Leave blank if using host and port
	// Socket string `json:"socket"`

	// The connection timeout to a ssdb server (seconds)
	Timeout int `json:"timeout"`

	// Maximum number of connections
	MaxConn int `json:"maxconn"`
}
