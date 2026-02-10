[Setup]
AppName=MirrorBox
AppVersion=1.0.0
AppPublisher=excellgene
DefaultDirName={autopf}\MirrorBox
DefaultGroupName=MirrorBox
UninstallDisplayIcon={app}\mirrorBox.exe
SourceDir=..\..\src\cmd\mirrorbox
OutputBaseFilename=mirrorBox-setup
OutputDir=.
SetupIconFile=winres\icon.ico
Compression=lzma2
SolidCompression=yes
ArchitecturesInstallIn64BitMode=x64compatible

[Files]
Source: "mirrorBox.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\MirrorBox"; Filename: "{app}\mirrorBox.exe"
Name: "{autodesktop}\MirrorBox"; Filename: "{app}\mirrorBox.exe"

[Run]
Filename: "{app}\mirrorBox.exe"; Description: "Launch MirrorBox"; Flags: nowait postinstall skipifsilent
