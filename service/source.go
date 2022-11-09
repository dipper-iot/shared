package service

import "os"

var (
	mapService = make(map[string]string)
)

const (
	ApiGatewayService    = "ApiGatewayService"
	AuthService          = "AuthService"
	RoleService          = "RoleService"
	UserService          = "UserService"
	DeviceService        = "DeviceService"
	ModelLogService      = "ModelLogService"
	TelemetryLogService  = "TelemetryLogService"
	MailService          = "MailService"
	NamespaceService     = "NamespaceService"
	ModelService         = "ModelService"
	UiService            = "UiService"
	AlertService         = "AlertService"
	OTPService           = "OTPService"
	MQTTTransportService = "MQTTTransportService"
	QueueService         = "QueueService"
	BrokerService        = "BrokerService"
)

func init() {
	setHost(AuthService, "AUTH_SERVICE_GRPC", "localhost:19000")
	setHost(UserService, "USER_SERVICE_GRPC", "localhost:19001")
	setHost(RoleService, "ROLE_SERVICE_GRPC", "localhost:19002")
	setHost(DeviceService, "DEVICE_SERVICE_GRPC", "localhost:19003")
	setHost(MailService, "MAIL_SERVICE_GRPC", "localhost:19004")
	setHost(NamespaceService, "NAMESPACE_SERVICE_GRPC", "localhost:19005")
	setHost(ModelService, "MODEL_SERVICE_GRPC", "localhost:19006")
	setHost(ModelLogService, "MODEL_LOG_SERVICE_GRPC", "localhost:19007")
	setHost(ApiGatewayService, "API_GATEWAY_SERVICE_GRPC", "localhost:18080")
	setHost(TelemetryLogService, "TELEMETRY_LOG_SERVICE_GRPC", "localhost:19008")
	setHost(OTPService, "OTP_SERVICE_GRPC", "localhost:19009")
	setHost(UiService, "UI_SERVICE_GRPC", "localhost:19010")
	setHost(AlertService, "ALERT_SERVICE_GRPC", "localhost:19011")
	setHost(MQTTTransportService, "MQTT_TRANSPORT_SERVICE_GRPC", "localhost:19012")
	setHost(QueueService, "QUEUE_SERVICE_GRPC", "localhost:19013")
	setHost(BrokerService, "BROKER_SERVICE_GRPC", "localhost:19014")
}

func setHost(serviceName string, env string, defaultAddress string) {
	host, success := os.LookupEnv(env)
	if !success {
		host = defaultAddress
	}
	mapService[serviceName] = host
}

func AddSource(name string, address string) {
	mapService[name] = address
}
