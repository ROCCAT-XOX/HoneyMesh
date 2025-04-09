package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"math/rand"
	"time"
)

// RealisticSensorData erzeugt realistischere Testdaten für die Bienenstöcke
func insertRealisticSensorData(client *mongo.Client, days int) error {
	collection := client.Database("HoneyMesh").Collection("data")

	// Konstanten für realistische Werte
	const (
		baseWeight     = 60.0 // Basis-Gewicht eines Bienenstocks in kg
		dailyVariation = 0.3  // Tägliche Variation (z.B. durch Nektar/Honig)
		tempVariation  = 5.0  // Temperaturschwankung im Bienenstock
		humidityBase   = 65.0 // Durchschnittliche Luftfeuchtigkeit im Bienenstock
	)

	// Aktuelle Zeit für die Zeitstempel
	now := time.Now()

	// Generiere zufällige Startwerte für jeden Bienenstock
	nodeCount := 5
	nodeWeights := make([]float64, nodeCount)
	for i := 0; i < nodeCount; i++ {
		// Jeder Bienenstock hat ein leicht unterschiedliches Startgewicht
		nodeWeights[i] = baseWeight + rand.Float64()*10.0
	}

	// Saisonale Faktoren (z.B. Frühling, Sommer haben mehr Nektar)
	seasonalFactor := func(date time.Time) float64 {
		// Monat als Float zwischen 0-12
		month := float64(date.Month())

		// Saisonaler Faktor: höher im Frühling/Sommer (Monate 3-8), niedriger im Herbst/Winter
		return 1.0 + 0.5*math.Sin((month-3.0)*math.Pi/6.0)
	}

	// Temperaturmuster basierend auf Tageszeit und Jahreszeit
	tempPattern := func(date time.Time, isInternal bool) float64 {
		hour := float64(date.Hour())
		month := float64(date.Month())

		// Tägliches Muster: wärmer am Tag, kühler in der Nacht
		dailyPattern := math.Sin((hour - 6.0) * math.Pi / 12.0) // Maximum um 12 Uhr

		// Jahreszeitliches Muster: wärmer im Sommer, kälter im Winter
		yearlyPattern := math.Sin((month - 3.0) * math.Pi / 6.0) // Maximum im Juli

		if isInternal {
			// Innentemperatur: bleibt relativ stabil in einem Bienenstock
			baseTemp := 34.0 // Bienen halten den Stock um 34-35°C
			return baseTemp + dailyPattern*1.0 + yearlyPattern*2.0
		} else {
			// Außentemperatur: folgt stärker dem Tages- und Jahreszyklus
			baseTemp := 15.0
			return baseTemp + dailyPattern*8.0 + yearlyPattern*15.0
		}
	}

	// Luftfeuchtigkeitsmuster
	humidityPattern := func(date time.Time) float64 {
		// Feuchtigkeit variiert mit Temperatur, Tageszeit und Wetter
		hour := float64(date.Hour())
		base := humidityBase

		// Feuchter in der Nacht, trockener am Tag
		dailyVariation := -10.0 * math.Sin((hour-6.0)*math.Pi/12.0)

		// Zufällige Schwankungen für wetterbedingte Änderungen
		randomVariation := rand.Float64()*10.0 - 5.0

		humid := base + dailyVariation + randomVariation
		// Begrenze auf realistischen Bereich
		return math.Max(30.0, math.Min(95.0, humid))
	}

	// WiFi-Signalstärke simulieren
	wifiStrengthPattern := func() float64 {
		// Werte zwischen -30 (sehr gut) und -90 (sehr schwach)
		baseStrength := -50.0 - rand.Float64()*20.0

		// Zufällige Schwankungen
		variation := rand.Float64()*10.0 - 5.0

		return baseStrength + variation
	}

	// Generiere Datenpunkte für jeden Tag mit Stunden-Intervallen
	for day := 0; day < days; day++ {
		// Rückwärts von heute
		currentDate := now.AddDate(0, 0, -day)

		// Wetterbedingte Variation für diesen Tag
		weatherFactor := rand.Float64()*0.4 + 0.8 // 0.8 bis 1.2

		// Honigzuwachs an produktiven Tagen
		honeyGrowthToday := rand.Float64() > 0.7 // 30% Chancen für ein produktiven Tag

		for hour := 0; hour < 24; hour += 3 { // Alle 3 Stunden
			// Aktuelle Zeit für diesen Datenpunkt
			dataTime := currentDate.Add(time.Duration(-hour) * time.Hour)
			currentDate := dataTime.Format("2006-01-02")
			currentTime := dataTime.Format("15:04:05") + fmt.Sprintf("%09d", time.Now().UnixNano()%1e9)

			// Saisonaler Faktor für diesen Tag
			seasonFactor := seasonalFactor(dataTime)

			var nodes []NodeData
			for i := 0; i < nodeCount; i++ {
				// Leichte Veränderung des Gewichts basierend auf Tageszeit/Saison
				dailyChange := rand.Float64()*dailyVariation*2.0 - dailyVariation

				// Füge Honig hinzu, wenn es ein produktiver Tag ist
				honeyBonus := 0.0
				if honeyGrowthToday && hour >= 8 && hour <= 18 {
					honeyBonus = 0.05 * seasonFactor // Mehr Honig in produktiven Monaten
				}

				// Aktualisiere das Gewicht für diesen Bienenstock
				nodeWeights[i] += dailyChange + honeyBonus

				// Sensorwerte für diesen Bienenstock
				tempIn := tempPattern(dataTime, true) + rand.Float64()*tempVariation/5.0
				tempOut := tempPattern(dataTime, false) + rand.Float64()*tempVariation
				humidity := humidityPattern(dataTime) * weatherFactor
				airQuality := 100.0 + rand.Float64()*50.0 // Einfacher Wert für Luftqualität

				// Wifi-Signalstärke - nicht abhängig von Wetter/Zeit
				wifiStrength := wifiStrengthPattern()

				// Node-Daten erstellen
				node := NodeData{
					NodeID:       fmt.Sprintf("%d", i+1),
					Weight:       nodeWeights[i],
					TempOut:      tempOut,
					TempIn:       tempIn,
					Humidity:     humidity,
					AirQuality:   airQuality,
					WifiStrength: wifiStrength,
					Date:         currentDate,
					Time:         fmt.Sprintf("%sT%s", currentDate, currentTime[:8]),
				}

				nodes = append(nodes, node)
			}

			// Sensor-Daten für alle Bienenstöcke zu diesem Zeitpunkt
			sensorData := SensorData{
				Date:  currentDate,
				Time:  currentTime[:8],
				Nodes: nodes,
			}

			// In die Datenbank einfügen
			_, err := collection.InsertOne(context.Background(), sensorData)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Ersetze die bestehende insertSampleData Funktion durch diese
func insertSampleData(client *mongo.Client) error {
	// Generiere Daten für die letzten 90 Tage
	return insertRealisticSensorData(client, 90)
}

// Abfrage 24h
func getSensorDataFromLast24Hours(client *mongo.Client) ([]SensorData, error) {
	collection := client.Database("HoneyMesh").Collection("data")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Berechnen der Daten für die letzten 24 Stunden
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	// Filter für Dokumente des aktuellen und des vorherigen Tages
	filter := bson.M{
		"date": bson.M{
			"$gte": yesterday.Format("2006-01-02"),
		},
	}
	options := options.Find().SetSort(bson.D{{"date", 1}, {"time", 1}})

	var results []SensorData
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var sensorData SensorData
		if err := cursor.Decode(&sensorData); err != nil {
			return nil, err
		}
		// Programmatische Filterung der Daten, um sicherzustellen, dass sie innerhalb der letzten 24 Stunden liegen
		sensorDateTimeStr := sensorData.Date + "T" + sensorData.Time
		sensorDateTime, err := time.Parse("2006-01-02T15:04:05", sensorDateTimeStr)
		if err != nil {
			log.Printf("Fehler beim Parsen des Datums/der Zeit: %v", err)
			continue
		}
		if sensorDateTime.After(yesterday) {
			results = append(results, sensorData)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	fmt.Println(results)
	return results, nil
}

// Dynamische Abfrage, Abfrage durch Stunden
func getFilteredSensorData(client *mongo.Client, hours int) ([]SensorData, error) {
	collection := client.Database("HoneyMesh").Collection("data")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Berechnen der Zeit für die letzten X Stunden, basierend auf dem Parameter
	targetTime := time.Now().Add(time.Duration(-hours) * time.Hour)

	filter := bson.M{
		"date": bson.M{
			"$gte": targetTime.Format("2006-01-02"),
		},
	}
	options := options.Find().SetSort(bson.D{{"date", 1}, {"time", 1}})

	var weights []SensorData
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var sensorData SensorData
		if err := cursor.Decode(&sensorData); err != nil {
			return nil, err
		}
		sensorDateTimeStr := sensorData.Date + "T" + sensorData.Time
		sensorDateTime, err := time.Parse("2006-01-02T15:04:05", sensorDateTimeStr)
		if err != nil {
			log.Printf("Fehler beim Parsen des Datums/der Zeit: %v", err)
			continue
		}
		if sensorDateTime.After(targetTime) {
			weights = append(weights, sensorData)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return weights, nil
}

func getLatestWeightForEachNode(client *mongo.Client) ([]bson.M, error) {
	collection := client.Database("HoneyMesh").Collection("data")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Aggregations-Pipeline
	pipeline := mongo.Pipeline{
		// Schritt 1: Dekonstruiere das Array 'nodes', um mit einzelnen Nodes arbeiten zu können
		{{"$unwind", "$nodes"}},
		// Schritt 2: Sortiere die Dokumente absteigend nach Datum und Zeit, um den neuesten Eintrag zuerst zu haben
		{{"$sort", bson.D{{"date", -1}, {"time", -1}}}},
		// Schritt 3: Gruppiere die Dokumente nach 'nodeid', und nimm den ersten (also neuesten) Eintrag für jedes Node
		{{"$group", bson.M{
			"_id":          "$nodes.nodeid",
			"latestWeight": bson.M{"$first": "$nodes.weight"},
			"latestDate":   bson.M{"$first": "$date"},
			"latestTime":   bson.M{"$first": "$time"},
		}}},
		// Optional: Sortiere die Ergebnisse nach NodeID in aufsteigender Reihenfolge, falls benötigt
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	var results []bson.M
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
