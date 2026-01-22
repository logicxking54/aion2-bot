#include <Keyboard.h>
#include <Mouse.h>

String buf = "";
bool killed = false;

void setup() {
  Serial.begin(115200);
  Keyboard.begin();
  Mouse.begin();
}

void loop() {
  while (Serial.available()) {
    char c = Serial.read();

    if (c == '\n') {
      handle(buf);
      buf = "";
    } else {
      buf += c;
    }
  }
}

void handle(String cmd) {
  if (cmd == "KILL") {
    emergencyStop();
    killed = true;
    return;
  }

  if (killed) return;

  if (cmd.length() < 3) return;

  char type = cmd.charAt(0);

  if (type == 'K') {
    char key = cmd.charAt(2);
    Keyboard.press(key);
    Keyboard.release(key);
  }

  else if (type == 'D') {
    char key = cmd.charAt(2);
    Keyboard.press(key);
  }

  else if (type == 'U') {
    char key = cmd.charAt(2);
    Keyboard.release(key);
  }

  else if (type == 'C') {
    int p = cmd.indexOf('|', 2);
    String mod = cmd.substring(2, p);
    char key = cmd.charAt(p + 1);

    if (mod == "CTRL")  Keyboard.press(KEY_LEFT_CTRL);
    if (mod == "SHIFT") Keyboard.press(KEY_LEFT_SHIFT);
    if (mod == "ALT")   Keyboard.press(KEY_LEFT_ALT);
    if (mod == "WIN")   Keyboard.press(KEY_LEFT_GUI);

    Keyboard.press(key);
    Keyboard.release(key);
    Keyboard.releaseAll();
  }

  else if (type == 'M') {
    int p = cmd.indexOf('|', 2);
    int dx = cmd.substring(2, p).toInt();
    int dy = cmd.substring(p + 1).toInt();
    Mouse.move(dx, dy);
  }
}

void emergencyStop() {
  Keyboard.releaseAll();

  Mouse.release(MOUSE_LEFT);
  Mouse.release(MOUSE_RIGHT);
  Mouse.release(MOUSE_MIDDLE);

  Mouse.move(0, 0);
}