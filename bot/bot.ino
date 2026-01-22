#include <Keyboard.h>

void setup() {
  Serial.begin(115200);
  Keyboard.begin();
}

void loop() {
  if (Serial.available()) {
    char cmd = Serial.read();

    if (cmd == 'W') {
      humanPressW();
    }
  }
}

void humanPressW() {
  delay(random(30, 60));          // hover แบบคน
  Keyboard.press('w');            // กด W ค้าง
  delay(random(120, 200));        // hold (สำคัญ)
  Keyboard.release('w');          // ปล่อย
  delay(random(50, 100));         // after delay
}
