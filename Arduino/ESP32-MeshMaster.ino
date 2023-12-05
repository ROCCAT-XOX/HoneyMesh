#include <painlessMesh.h>
#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>

#define MESH_PREFIX "meshNetwork"
#define MESH_PASSWORD "meshPassword"
#define MESH_PORT 5555

#define OLED_RESET -1
Adafruit_SSD1306 display(128, 32, &Wire, OLED_RESET);

painlessMesh mesh;

void receivedCallback(const uint32_t &from, const String &msg) {
  Serial.printf("Received message from Node %u: %s\n", from, msg.c_str());

  // Nachricht im vorgegebenen Format: fixedID,weight,tempIn,tempOut
  // Beispiel: "1,100,25,20"
  Serial2.println(msg); // Sendet die empfangene Nachricht an den Webserver
}

void newConnectionCallback(const uint32_t &id) {
  Serial.printf("New Connection, nodeId = %u\n", id);
  updateDisplay();
}

void changedConnectionCallback() {
  Serial.printf("Changed connections\n");
  updateDisplay();
}

void nodeTimeAdjustedCallback(const int32_t &offset) {
  Serial.printf("Adjusted time %u. Offset = %d\n", mesh.getNodeTime(), offset);
}

void updateDisplay() {
  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(WHITE);
  display.setCursor(0,0);
  display.println("Mesh Nodes:");
  
  auto nodeList = mesh.getNodeList();
  display.println(nodeList.size());

  int y = 16; // Startposition für die erste ID
  for(auto &nodeId : nodeList) {
    if(y > 31) break; // Begrenzung der Anzeige auf den Displaybereich
    display.setCursor(0, y);
    display.print("ID: ");
    display.println(nodeId);
    y += 8; // Verschiebt die Y-Position für die nächste ID
  }

  display.display();
}

void setup() {
  Serial.begin(115200);
  Serial2.begin(115200, SERIAL_8N1, 16, 17); // Verwendet GPIO 16 als TX und GPIO 17 als RX

  mesh.setDebugMsgTypes(ERROR | DEBUG); 
  mesh.init(MESH_PREFIX, MESH_PASSWORD, MESH_PORT);
  mesh.onReceive(&receivedCallback);
  mesh.onNewConnection(&newConnectionCallback);
  mesh.onChangedConnections(&changedConnectionCallback);
  mesh.onNodeTimeAdjusted(&nodeTimeAdjustedCallback);

  Wire.begin(21, 22);
  if(!display.begin(SSD1306_SWITCHCAPVCC, 0x3C)) {
    Serial.println("SSD1306 allocation failed");
    for(;;);
  }
  updateDisplay();
}

void loop() {
  mesh.update();
}
