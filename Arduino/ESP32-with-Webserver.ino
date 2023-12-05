#include <WiFi.h>
#include <WebServer.h>
#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>
#include <time.h>
#include <map>
#include <ArduinoJson.h>  // Stellen Sie sicher, dass diese Bibliothek eingebunden ist

#define SCREEN_WIDTH 128
#define SCREEN_HEIGHT 32
#define OLED_RESET -1
Adafruit_SSD1306 display(SCREEN_WIDTH, SCREEN_HEIGHT, &Wire, OLED_RESET);

const char* ssid = "Rueff-Village";
const char* password = "Pedro2020!";
WebServer server(80);

const char* ntpServer = "pool.ntp.org";
const long gmtOffset_sec = 3600;
const int daylightOffset_sec = 3600;

#define RX_PIN 16
#define TX_PIN 17

struct NodeData {
  int weight;
  int tempIn;
  int tempOut;
};

std::map<String, NodeData> nodeDataMap;

void setup() {
  Serial.begin(115200);
  Serial2.begin(115200, SERIAL_8N1, RX_PIN, TX_PIN);

  Wire.begin(21, 22);
  
  if(!display.begin(SSD1306_SWITCHCAPVCC, 0x3C)) { 
    Serial.println(F("SSD1306 allocation failed"));
    for(;;);
  }
  display.display();
  delay(2000);

  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.println("Connecting to WiFi...");
  }
  Serial.println("Connected to WiFi");

  configTime(gmtOffset_sec, daylightOffset_sec, ntpServer);

  server.on("/", HTTP_GET, handleRoot);
  server.begin();
  Serial.println("HTTP server started");

  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(WHITE);
  display.setCursor(0,0);
  display.println("IP Address:");
  display.println(WiFi.localIP());
  display.println("Webserver started");
  display.display();
}

void handleRoot() {
  if (!server.hasArg("format") || server.arg("format") == "json") {
    server.send(200, "application/json", createJSONData());
  } else if (server.arg("format") == "text") {
    String dataFromMeshMaster = readSerialData();
    server.send(200, "text/plain", dataFromMeshMaster);
  } else {
    server.send(404, "text/plain", "Invalid format");
  }
}

String createJSONData() {
  struct tm timeinfo;
  if (!getLocalTime(&timeinfo)) {
    return "{\"error\":\"Failed to obtain time\"}";
  }
  char timeString[64];
  char dateString[64];
  strftime(timeString, sizeof(timeString), "%H:%M:%S", &timeinfo);
  strftime(dateString, sizeof(dateString), "%Y-%m-%d", &timeinfo);

  DynamicJsonDocument doc(1024);
  doc["date"] = dateString;
  doc["time"] = timeString;

  for (auto& pair : nodeDataMap) {
    JsonObject node = doc.createNestedObject("Node" + pair.first);
    node["weight"] = pair.second.weight;
    node["tempIn"] = pair.second.tempIn;
    node["tempOut"] = pair.second.tempOut;
  }

  String json;
  serializeJson(doc, json);
  return json;
}

String readSerialData() {
  String data = "";
  while (Serial2.available()) {
    String line = Serial2.readStringUntil('\n');
    if (line.length() > 0) {
      int firstComma = line.indexOf(',');
      int secondComma = line.indexOf(',', firstComma + 1);
      int thirdComma = line.indexOf(',', secondComma + 1);

      String nodeID = line.substring(0, firstComma);
      NodeData nd;
      nd.weight = line.substring(firstComma + 1, secondComma).toInt();
      nd.tempIn = line.substring(secondComma + 1, thirdComma).toInt();
      nd.tempOut = line.substring(thirdComma + 1).toInt();

      nodeDataMap[nodeID] = nd;
    }
  }
  return data;
}

void loop() {
  readSerialData();
  server.handleClient();
}
