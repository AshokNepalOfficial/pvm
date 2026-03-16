#define MyAppName "PVM for Windows"
#define MyAppShortName "pvm"
#define MyAppLCShortName "pvm"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "Eskycode Technology Pvt. Ltd"
#define MyAppURL "https://github.com/ashoknepalofficial/pvm"
#define MyAppExeName "pvm.exe"
#define MyIcon "bin\pvm.ico"
#define MyAppId "ccfb402b-a5b4-4dd3-a45e-06e832bd9bee"
#define ProjectRoot "C:\Users\eskyc\Downloads\pvm-master\pvm-master"

[Setup]
PrivilegesRequired=admin
AppId={#MyAppId}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppCopyright=Copyright (C) 2026 {#MyAppPublisher}
AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={localappdata}\{#MyAppShortName}
DisableDirPage=no
DefaultGroupName={#MyAppName}
AllowNoIcons=yes
LicenseFile={#ProjectRoot}\LICENSE.txt
OutputDir={#ProjectRoot}\dist\{#MyAppVersion}
OutputBaseFilename=pvm-setup
SetupIconFile={#ProjectRoot}\{#MyIcon}
Compression=lzma
SolidCompression=yes
ChangesEnvironment=yes
DisableProgramGroupPage=yes
ArchitecturesInstallIn64BitMode=x64compatible
UninstallDisplayIcon={app}\{#MyIcon}

; Version information
VersionInfoVersion={#MyAppVersion}.0
VersionInfoCompany={#MyAppPublisher}
VersionInfoDescription=PHP version manager for Windows
VersionInfoProductName={#MyAppName}
VersionInfoProductTextVersion={#MyAppVersion}
VersionInfoOriginalFileName=pvm-setup.exe
VersionInfoCopyright=Copyright (C) 2026 {#MyAppPublisher}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "{#ProjectRoot}\bin\*"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs; Excludes: "{#ProjectRoot}\bin\install.cmd"

[Icons]
Name: "{group}\{#MyAppShortName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{#ProjectRoot}\{#MyIcon}"
Name: "{group}\Uninstall {#MyAppShortName}"; Filename: "{uninstallexe}"

[Registry]
Root: HKCR; Subkey: "{#MyAppShortName}"; ValueType: string; ValueName: ""; ValueData: "URL:pvm"; Flags: uninsdeletekey
Root: HKCR; Subkey: "{#MyAppShortName}"; ValueType: string; ValueName: "URL Protocol"; ValueData: ""; Flags: uninsdeletekey
Root: HKCR; Subkey: "{#MyAppShortName}\DefaultIcon"; ValueType: string; ValueName: ""; ValueData: "{app}\{#MyAppExeName},0"; Flags: uninsdeletekey
Root: HKCR; Subkey: "{#MyAppShortName}\shell\launch\command"; ValueType: string; ValueName: ""; ValueData: """{app}\{#MyAppExeName}"" ""%1"""; Flags: uninsdeletekey

[Code]
var
  SymlinkPage: TInputDirWizardPage;

function GetCurrentYear(Param: String): String;
begin
  result := GetDateTimeString('yyyy', '-', ':');
end;

function IsDirEmpty(dir: string): Boolean;


var
  FindRec: TFindRec;
  ct: Integer;
begin
  ct := 0;
  if FindFirst(ExpandConstant(dir + '\*'), FindRec) then
  try
    repeat
      if FindRec.Attributes and FILE_ATTRIBUTE_DIRECTORY = 0 then
        ct := ct+1;
    until
      not FindNext(FindRec);
  finally
    FindClose(FindRec);
    Result := ct = 0;
  end;
end;

procedure InitializeWizard;
begin
  { Symlink Page – Active PHP Version Location }
  SymlinkPage := CreateInputDirPage(wpSelectDir,
    'Active PHP Version Location',
    'The active PHP version will always be available here.',
    'Select the folder in which Setup should create the symlink, then click Next.',
    False, '');
  SymlinkPage.Add('This directory will automatically be added to your system PATH.');
  SymlinkPage.Values[0] := ExpandConstant('C:\pvm\php');
end;

function NextButtonClick(CurPageID: Integer): Boolean;
begin
  Result := True; // allow Next by default

  if CurPageID = SymlinkPage.ID then
  begin
    if DirExists(SymlinkPage.Values[0]) and not IsDirEmpty(SymlinkPage.Values[0]) then
    begin
      MsgBox('The selected folder is not empty. Please choose another path.', mbError, MB_OK);
      Result := False;
    end;
  end;
end;

function getSymLink(Param: string): string;
begin
  Result := SymlinkPage.Values[0];
end;

procedure CurStepChanged(CurStep: TSetupStep);
var
  path: string;
begin
  if CurStep = ssPostInstall then
  begin
    SaveStringToFile(ExpandConstant('{app}\settings.txt'),
      'root: ' + ExpandConstant('{app}') + #13#10 +
      'path: ' + SymlinkPage.Values[0] + #13#10,
      False);

    { Write environment variables }
    RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'PVM_HOME', ExpandConstant('{app}'));
    RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'PVM_SYMLINK', SymlinkPage.Values[0]);

    { Update PATH }
    RegQueryStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', path);
    if Pos('%PVM_HOME%',path)=0 then path := path+';%PVM_HOME%';
    if Pos('%PVM_SYMLINK%',path)=0 then path := path+';%PVM_SYMLINK%';
    RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', path);
  end;
end;

[Run]
Filename: "{cmd}"; Parameters: "/C mklink /D ""{code:getSymLink}"" ""{app}\bin"""; Flags: runhidden waituntilterminated
Filename: "powershell.exe"; Parameters: "-NoExit -Command refreshenv; cls; Write-Host 'Welcome to PVM for Windows'"; Description: "Open with Powershell"; Flags: postinstall skipifsilent

[UninstallDelete]
Type: filesandordirs; Name: "{app}";
Type: filesandordirs; Name: "{localappdata}\.pvm";