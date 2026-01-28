# SambaSync

Cross platform tool to sync folders periodically


## Build & Run

1. Run the application

    Windows
    
        go run .\src\cmd\app\main.go

    MacOS

        go run ./src/cmd/app/main.go

1. Build a release version

    Windows

        cd src/cmd/app

        fyne package -name sambaSync -os windows -icon ../../internal/ui/icons/appicon.png 

    MacOS

        cd src/cmd/app

        fyne package -name sambaSync -os darwin -icon ../../internal/ui/icons/appicon.png --metadata LSUIElement=true
