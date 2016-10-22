touchpad-cover
--------------

Disable your touchpad while typing. Linux only.

You probably want `syndaemon` instead of this. But if that doesn't work for you (as it didn't for me), this program might.


## Installation

If you have Go installed, you can use:

    go get -u github.com/jeremyschlatter/touchpad-cover

Or [download a binary](https://github.com/jeremyschlatter/touchpad-cover/releases) from GitHub.

## Configuration

You need to identify the names of two things in your system:

**The device name of your touchpad**

Use `xinput list` to find it. My input list looks like this:

```
⎡ Virtual core pointer                    	id=2	[master pointer  (3)]
⎜   ↳ Virtual core XTEST pointer              	id=4	[slave  pointer  (2)]
⎜   ↳ Atmel maXTouch Touchscreen              	id=12	[slave  pointer  (2)]
⎜   ↳ Atmel maXTouch Touchpad                 	id=11	[slave  pointer  (2)]
⎣ Virtual core keyboard                   	id=3	[master keyboard (2)]
    ↳ Virtual core XTEST keyboard             	id=5	[slave  keyboard (3)]
    ↳ Power Button                            	id=6	[slave  keyboard (3)]
    ↳ Video Bus                               	id=7	[slave  keyboard (3)]
    ↳ Power Button                            	id=8	[slave  keyboard (3)]
    ↳ Sleep Button                            	id=9	[slave  keyboard (3)]
    ↳ Sleep Button                            	id=10	[slave  keyboard (3)]
    ↳ Chromebook HD WebCam                    	id=13	[slave  keyboard (3)]
    ↳ AT Translated Set 2 keyboard            	id=14	[slave  keyboard (3)]
```
Here, my touchpad name is "Atmel maXTouch Touchpad"

**The input event device file for your keyboard**

This should be located somewhere under `/dev/input`.

It might be `/dev/input/eventX`, where `X` is the id of your keyboard in `xinput list`. (So for the example above, this would be `/dev/input/event5`). You can also look for a symbolic name that looks like it includes "keyboard" in `/dev/input/by-path` or `/dev/input/by-id`. For example, on my system `/dev/input/by-path/platform-i8042-serio-0-event-kbd` is the file I want (notice the "kbd" in the name). To check whether a particular file is what you want, run `sudo cat <file>` and see if it outputs anything as you press keys. Or you can run `sudo touchpad-cover --touchpad-name <name from step 1> --dev-input-keyboard <event file to try> --verbose` and again see if it outputs anything when you press keys.

## Use

When you have identified both of the above, run:

    sudo touchpad-cover --touchpad-name <touchpad name> --dev-input-keyboard <key event file>
    
If it works as expected, set up that command to run on startup and you should be good to go.
