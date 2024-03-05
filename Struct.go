package main

type SensorData struct {
	Date  string     `bson:"date"`
	Time  string     `bson:"time"`
	Nodes []NodeData `bson:"nodes"`
}

type NodeData struct {
	NodeID       string  `bson:"nodeid"`
	Weight       float64 `bson:"weight"`
	TempOut      float64 `bson:"tempout"`
	TempIn       float64 `bson:"tempin"`
	Humidity     float64 `bson:"humidity"`
	AirQuality   float64 `bson:"airquality"`
	WifiStrength float64 `bson:"wifistrength"`
	Date         string  `bson:"date"`
	Time         string  `bson:"time"`
}

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"` // Hinweis: In einer Produktionsumgebung sollte das Passwort gehasht gespeichert werden
}
