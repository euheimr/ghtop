# ghtop (go htop)
A terminal-based activity monitor inspired by [gotop](https://github.com/xxxserxxx/gotop), [htop](https://hisham.hm/htop/), [gtop](https://github.com/aksakalli/gtop) and [vtop](https://github.com/MrRio/vtop), entirely written in [Go](https://golang.org/).

---

**Also, a special _THANK YOU_ to [cjbassi](https://github.com/cjbassi/) for making the original [gotop](https://github.com/cjbassi/gotop) and [xxxserxxx](https://github.com/xxxserxxx/) for maintaining [gotop](https://github.com/xxxserxxx/gotop)!**

This application is inspired by `gotop`, but I wanted to make my own version with improvements (sysinfo, bars instead of mostly graphs, and code structure changes) and by using the newer Text UI (TUI) go package [tview](https://github.com/rivo/tview)!

Subsequently, I had to draw a lot of my own widgets/primitives that would _otherwise_ be included with the TUI go package [termui](https://github.com/gizak/termui) used by [gotop](https://github.com/xxxserxxx/gotop).

## Why use this over gotop?

!TODO: do performance testing between gotop and ghtop

1. ?performance
2. Works for Windows, Linux and macOS (explicitly tested for all 3)
3. This is written for Go version [1.20](https://go.dev/dl/)+, and thus takes advantage of the newer features of Go
4. Most of the `gotop` dependencies are fairly out-of-date (uses old and soft-deprecated Go code)
5. Uses the newer `tview` text UI go package (I personally like the general features and API it provides over `termui`)

**Most importantly, I wanted to learn Go.** 

This project helped me a lot and pushed me to learn most of the features of Go.

## Install
If you're on windows, the new windows terminal is packaged with this application under the `terminal` directory.
This is executed in the end to start `ghtop`.

**Please note: Unicode (Braille characters) isn't supported on the old Windows console (`cmd` / `conhost.exe`)!**


1. !todo: Download & Install from [Releases]()

  **OR:** 

2. Build it yourself

   a. Windows
      
      1. `cmd` or `conhost.exe`:
   
       build.bat && install.bat

      2. `Powershell`:
       build.bat; install.bat

   b. Linux & MacOS

       build.sh

3. Install & Add to path

       install.sh


## Built With:

 - tcell
 - tview
 - gopsutil
 - drawille-go


## Reference

### Design Mockup

<div>
<img src="./docs/ghtop.png" alt="design mockup"/>
</div>