# v2.0.9

* Fixed bug. 

# v2.0.8

* Added 'fast-diagonals' #26

# v2.0.7

* Ignored capitalisation.

# v2.0.5

* New tray icon painted by **Koka (gavrilovne@gmail.com)**

# v2.0.4

* Installer will automatically kill Mouseable process.
* Fixed a bug where Attach not working.

# v2.0.3

* Added data **(cursor|wheel|teleport)-factor**. You can set horizontal/vertical speed differently like
  **TeleportDistanceH/V** in previous version. #19


* Now accelerations not carried over. #17

# v2.0.2

* Fixed high CPU usage.

# v2.0.1

### UI

Dropped sloppy UI, embed [monaco-editor](https://microsoft.github.io/monaco-editor). Now you can edit json file directly
with comfortable suggestion.

### New Logic

All logic have been newly created.

To reduce cpu usage and cursor accuracy, **friction was deleted.**  There are **acceleration and max-speed**
now. It could be uncomfortable, I'll get your opinion and revise it.

### Abstraction

It has become a more suitable code for maintenance.

### Etc.

There are many changes that have not been mentioned. There may be a bug, please let me know. I will fix it.

# v1.0.14

### Bug Fixes

* Fixed a bug where window not open.
* Added option to turn off overlay.
* Changed overlay behavior to always follow cursor.

### New Features

* Added Terminate button to GUI.

# v1.0.13

### Bug Fixes

* Fixed a bug where mouseable not launched.

# v1.0.12

### Bug Fixes

* Fixed a bug where horizontal scrolling was not working.

# v1.0.10

### Brand New UI

Mouseable now use [lorca](https://github.com/zserge/lorca) to draw GUI.

### New Features

* Added acceleration and friction to wheel.
* Added 8 way Teleport.
* Added horizontal and vertical data.
* Added double press support.
* Added Win + Ctrl + Alt + Shift combination.

### Feature Updates

* Rename Flash to Teleport.

# v1.0.9

### New Features

* You can no longer run multiple Mouseable at the same time.

### Feature Updates

* Removed 'Save' button. Key changes now take effect immediately.

### Bug Fixes

* Fixed a bug where hotkeys could not be set.
* Fixed a bug where the window get cut off.
* Fixed a bug where version was not displayed.

# v1.0.8

* Added GUI for keymap.
* Fixed crashing issue.
* Fixed bugs.

# v1.0.7

* Changed 'Flash' logic to jump fixed distance regardless of velocity.

# v1.0.6

* Added an overlay near the mouse when activated.
* Added 'Flash'.
* Fixed bugs.

# v1.0.5

* Rewrote loop logic to ensure a constant period.
* Changed default config.
* Fixed bugs.

# v1.0.4

* Rewrote whole program to use keyboard hook only when necessary.

# v1.0.3

* Fixed bugs.

# v1.0.2

* Fixed bugs.

# v1.0.1

First version
