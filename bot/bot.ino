#include <Keyboard.h>
#include <Mouse.h>

void setup() {
  Serial.begin(115200);
  Keyboard.begin();
  Mouse.begin();
}

void loop() {
  char action = Serial.read();

  if (action == 'D') {
    if (Serial.available() < 1) return;

    char key = Serial.read();

    if (key == 'L') {
      Mouse.press(MOUSE_LEFT);
    } else {
      Keyboard.press(key);
    }
  } else if (action == 'U') {
    if (Serial.available() < 1) return;
    
    char key = Serial.read();

    if (key == 'L') {
      Mouse.release(MOUSE_LEFT);
    } else {
      Keyboard.release(key);
    }
  } else if (action == 'M') {
    if (Serial.available() < 2) return;

    int8_t dx = Serial.read();
    int8_t dy = Serial.read();

    Mouse.move(dx, dy);
  }
}
