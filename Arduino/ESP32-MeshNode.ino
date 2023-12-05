#include <painlessMesh.h>

#define MESH_PREFIX "meshNetwork"
#define MESH_PASSWORD "meshPassword"
#define MESH_PORT 5555

painlessMesh mesh;

unsigned long lastTime = 0;
unsigned long interval = 5000; // 5 Sekunden
int fixedID = 1; // Feste ID des Nodes

void setup() {
  Serial.begin(115200);
  mesh.setDebugMsgTypes(ERROR | DEBUG);
  mesh.init(MESH_PREFIX, MESH_PASSWORD, MESH_PORT);
}

void loop() {
  mesh.update();

  if (millis() - lastTime > interval) {
    lastTime = millis();

    int weight = random(50, 121); // Gewicht zwischen 50 und 120
    int tempIn = random(0, 41);   // TempIn zwischen 0 und 40 Grad
    int tempOut = random(-10, 41); // TempOut zwischen -10 und 40 Grad

    // Nachricht zusammenstellen
    String msg = String(fixedID) + "," + String(weight) + "," + String(tempIn) + "," + String(tempOut);

    mesh.sendBroadcast(msg);
    Serial.println(msg); // Zu Debug-Zwecken
  }
}
