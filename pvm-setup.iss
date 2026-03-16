#define MyAppName "PVM for Windows"
#define MyAppShortName "pvm"
#define MyAppVersion "1.2.0"
#define MyAppPublisher "Eskycode Technology Pvt. Ltd"
#define MyAppURL "https://github.com/ashoknepalofficial/pvm"
#define MyAppExeName "pvm.exe"
#define MyAppIcon "pvm.ico"
#define MyAppId "ccfb402b-a5b4-4dd3-a45e-06e832bd9bee"
#define ProjectRoot "C:\Users\eskyc\Downloads\pvm-master\pvm-master"

[Setup]
PrivilegesRequired=admin
AppId={#MyAppId}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={localappdata}\{#MyAppShortName}
DisableDirPage=no
DefaultGroupName={#MyAppName}
AllowNoIcons=yes
; Ensure this file exists at the path below before compiling
LicenseFile={#ProjectRoot}\LICENSE.txt
OutputDir={#ProjectRoot}\dist\{#MyAppVersion}
OutputBaseFilename=pvm-setup
SetupIconFile={#ProjectRoot}\bin\pvm.ico
Compression=lzma
SolidCompression=yes
ChangesEnvironment=yes
DisableProgramGroupPage=yes
ArchitecturesInstallIn64BitMode=x64compatible
UninstallDisplayIcon={app}\{#MyAppIcon}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Icons]
Name: "{group}\{#MyAppShortName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{app}\{#MyAppIcon}"
Name: "{group}\Uninstall {#MyAppShortName}"; Filename: "{uninstallexe}"

[Code]
var
  SymlinkPage: TInputDirWizardPage;

// Improved Download function using PowerShell (Handles TLS 1.2+ automatically)
function DownloadFile(const URL, Dest: string): Boolean;
var
  ErrorCode: Integer;
  PowerShellCmd: string;
begin
  // Create the directory if it doesn't exist
  ForceDirectories(ExtractFilePath(Dest));
  
  // PowerShell command to download file
  PowerShellCmd := Format('-ExecutionPolicy Bypass -Command "[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; (New-Object System.Net.WebClient).DownloadFile(''%s'', ''%s'');"', [URL, Dest]);
  
  // Execute PowerShell hidden
  Result := Exec('powershell.exe', PowerShellCmd, '', SW_HIDE, ewWaitUntilTerminated, ErrorCode);
  
  if not Result or (ErrorCode <> 0) then
  begin
    Log('Download failed for: ' + URL + ' Error Code: ' + IntToStr(ErrorCode));
    Result := False;
  end else
  begin
    Result := True;
  end;
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
      if (FindRec.Name <> '.') and (FindRec.Name <> '..') then
        ct := ct + 1;
    until not FindNext(FindRec);
  finally
    FindClose(FindRec);
    Result := ct = 0;
  end;
end;

procedure InitializeWizard;
begin
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
  Result := True;
  if CurPageID = SymlinkPage.ID then
  begin
    if DirExists(SymlinkPage.Values[0]) and (not IsDirEmpty(SymlinkPage.Values[0])) then
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
  path, BinPath, IconPath: string;
begin
  if CurStep = ssInstall then
  begin
    BinPath := ExpandConstant('{app}\{#MyAppExeName}');
    IconPath := ExpandConstant('{app}\{#MyAppIcon}');

    // 1. Download PVM binary
    // Note: Ensure this URL points to the actual raw .exe, not the setup installer itself
    WizardForm.StatusLabel.Caption := 'Downloading PVM binary...';
    if not DownloadFile('https://github.com/AshokNepalOfficial/pvm/releases/download/v1.2.0/pvm-setup.exe', BinPath) then
    begin
      MsgBox('Failed to download pvm.exe. Please check your internet connection.', mbError, MB_OK);
    end;

    // 2. Download icon
    WizardForm.StatusLabel.Caption := 'Downloading application icon...';
    DownloadFile('https://raw.githubusercontent.com/ashoknepalofficial/pvm/main/bin/pvm.ico', IconPath);
  end;

  if CurStep = ssPostInstall then
  begin
    SaveStringToFile(
      ExpandConstant('{app}\settings.txt'),
      'root: ' + ExpandConstant('{app}') + #13#10 +
      'path: ' + SymlinkPage.Values[0] + #13#10,
      False
    );

    RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'PVM_HOME', ExpandConstant('{app}'));
    RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'PVM_SYMLINK', SymlinkPage.Values[0]);

    if RegQueryStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', path) then
    begin
      if Pos('%PVM_HOME%', path) = 0 then path := path + ';%PVM_HOME%';
      if Pos('%PVM_SYMLINK%', path) = 0 then path := path + ';%PVM_SYMLINK%';
      RegWriteExpandStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', path);
    end;
  end;
end;

[Run]
Filename: "{cmd}"; Parameters: "/C if not exist ""{code:getSymLink}"" mkdir ""{code:getSymLink}"""; Flags: runhidden waituntilterminated
Filename: "{cmd}"; Parameters: "/C rmdir ""{code:getSymLink}"" 2>nul"; Flags: runhidden waituntilterminated
Filename: "{cmd}"; Parameters: "/C mklink /D ""{code:getSymLink}"" ""{app}"""; Flags: runhidden waituntilterminated
Filename: "{cmd}"; Parameters: "/K echo Welcome to PVM for Windows && ""{app}\{#MyAppExeName}"" help"; Description: "Open PVM in Command Prompt"; Flags: postinstall skipifsilent

[UninstallDelete]
Type: filesandordirs; Name: "{app}";
Type: filesandordirs; Name: "{localappdata}\.pvm";