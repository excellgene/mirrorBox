# MirrorBox

Cross platform tool to sync folders periodically

<div align="center">
  <img src="./assets/1.png" alt="MirrorBox Logo" width="128" height="128"/>
</div> 

![Settings](./assets/2.jpg)
![General Settings](./assets/3.jpg)
![Add folder 1](./assets/4.jpg)
![Add folder 2](./assets/5.jpg)


## Build & Run

1. Run the application

       go run ./src/cmd/app/main.go

1. Build package

        cd src/

        fyne package -name mirrorBox -os darwin -icon ./Icon.png ./cmd/app
