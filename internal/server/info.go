package server

type serverInfo struct{}

func GetInfoSection() serverInfo {
	return serverInfo{}
}
