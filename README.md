<div align="center">
  <img src="https://github.com/excellgene/mirrorBox/blob/main/src/cmd/mirrorbox/Icon.png" alt="App logo" width="100" height="100"/>
</div>

# MirrorBox

Cross platform tool to sync folders periodically

<div align="center">
  <img src="https://github.com/excellgene/mirrorBox/blob/main/assets/1.jpg" alt="Screen shot 1" width="300" height="200"/>
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
