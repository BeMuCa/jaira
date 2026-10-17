# Installs the latest jaira release binary on Windows.
#
#   irm https://raw.githubusercontent.com/BeMuCa/jaira/master/scripts/install.ps1 | iex
#
# The Windows twin of install.sh: downloads the release zip for this
# architecture, verifies its checksum, and puts jaira.exe in
# %LOCALAPPDATA%\Programs\jaira (override with JAIRA_INSTALL_DIR), then adds
# that directory to the user PATH. No admin rights, nothing beyond Windows
# PowerShell 5.1. Run it again to update in place. Then it runs
# `jaira roles install --global` so the role skills and the cockpit are set up
# too; JAIRA_NO_ROLES=1 installs only the binary.
#
# Everything lives in functions and errors are thrown, never `exit`ed: under
# `irm | iex` the script runs inside the caller's session, and an exit would
# close their window.

function Get-JairaArch {
	# A 32-bit PowerShell on 64-bit Windows reports x86 in PROCESSOR_ARCHITECTURE
	# and the real machine in PROCESSOR_ARCHITEW6432.
	$arch = $env:PROCESSOR_ARCHITEW6432
	if (-not $arch) { $arch = $env:PROCESSOR_ARCHITECTURE }
	switch ($arch) {
		'AMD64' { return 'amd64' }
		'ARM64' { return 'arm64' }
		default { throw "unsupported architecture '$arch'" }
	}
}

# Test-JairaChecksum throws unless the archive's SHA-256 matches its line in
# checksums.txt, so a truncated or tampered download never lands on PATH.
function Test-JairaChecksum([string]$ArchivePath, [string]$ChecksumsPath) {
	$name = Split-Path -Leaf $ArchivePath
	$want = $null
	foreach ($line in Get-Content -LiteralPath $ChecksumsPath) {
		$fields = $line -split '\s+'
		if ($fields.Count -ge 2 -and $fields[1] -eq $name) { $want = $fields[0]; break }
	}
	if (-not $want) { throw "$name is not in checksums.txt" }
	# .NET rather than Get-FileHash: Windows PowerShell 5.1 started from a pwsh 7
	# session inherits its module path and cannot load Get-FileHash at all.
	$sha = [Security.Cryptography.SHA256]::Create()
	$stream = [IO.File]::OpenRead($ArchivePath)
	try { $got = [BitConverter]::ToString($sha.ComputeHash($stream)) -replace '-', '' }
	finally { $stream.Dispose(); $sha.Dispose() }
	if ($got -ne $want) { throw "checksum mismatch for $name" }
}

# Add-JairaUserPath puts $Dir on the user PATH unless it is already there, and
# reports whether it changed anything. The value is read and written raw
# through the registry: [Environment]::SetEnvironmentVariable would expand
# entries like %USERPROFILE%\bin into fixed paths for good.
function Add-JairaUserPath([string]$Dir) {
	$key = 'HKCU:\Environment'
	$raw = (Get-Item -LiteralPath $key).GetValue('Path', '', 'DoNotExpandEnvironmentNames')
	$norm = $Dir.TrimEnd('\')
	foreach ($entry in $raw -split ';') {
		$e = [Environment]::ExpandEnvironmentVariables($entry).TrimEnd('\')
		if ($e -and $e -eq $norm) { return $false }
	}
	$new = if ($raw) { $raw.TrimEnd(';') + ';' + $Dir } else { $Dir }
	Set-ItemProperty -LiteralPath $key -Name Path -Value $new -Type ExpandString
	return $true
}

# Writing the registry tells no one; this broadcast is what makes a newly
# opened terminal see the new PATH without signing out.
function Send-JairaEnvChange {
	try {
		if (-not ('JairaInstall.Env' -as [type])) {
			Add-Type -Namespace JairaInstall -Name Env -MemberDefinition @'
[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam, uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
'@
		}
		$result = [UIntPtr]::Zero
		[void][JairaInstall.Env]::SendMessageTimeout([IntPtr]0xffff, 0x1a, [UIntPtr]::Zero, 'Environment', 2, 5000, [ref]$result)
	} catch {
		Write-Host 'note: open a new terminal (or sign out and in) for the PATH change to apply'
	}
}

function Install-Jaira {
	$ErrorActionPreference = 'Stop'
	# Windows PowerShell 5.1 still offers TLS 1.0 first, which GitHub refuses,
	# and redraws its progress bar so often that a download crawls.
	[Net.ServicePointManager]::SecurityProtocol = [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12
	$ProgressPreference = 'SilentlyContinue'
	$interactive = [Environment]::UserInteractive

	$repo = 'BeMuCa/jaira'
	$installDir = $env:JAIRA_INSTALL_DIR
	if (-not $installDir) { $installDir = Join-Path $env:LOCALAPPDATA 'Programs\jaira' }
	$arch = Get-JairaArch

	try {
		$tag = (Invoke-RestMethod -UseBasicParsing "https://api.github.com/repos/$repo/releases/latest").tag_name
	} catch { $tag = $null }
	if (-not $tag) {
		throw "no release found - install from source instead: go install github.com/$repo/cmd/jaira@latest"
	}
	$version = $tag.TrimStart('v')

	$archive = "jaira_${version}_windows_${arch}.zip"
	$base = "https://github.com/$repo/releases/download/$tag"

	$tmp = Join-Path ([IO.Path]::GetTempPath()) ("jaira-install-" + [Guid]::NewGuid())
	New-Item -ItemType Directory -Path $tmp | Out-Null
	try {
		Write-Host "downloading jaira $tag (windows/$arch)..."
		# Only the archive is big enough for a progress bar, and only when a
		# person is watching.
		if ($interactive) { $ProgressPreference = 'Continue' }
		Invoke-WebRequest -UseBasicParsing "$base/$archive" -OutFile (Join-Path $tmp $archive)
		$ProgressPreference = 'SilentlyContinue'
		Invoke-WebRequest -UseBasicParsing "$base/checksums.txt" -OutFile (Join-Path $tmp 'checksums.txt')
		Test-JairaChecksum (Join-Path $tmp $archive) (Join-Path $tmp 'checksums.txt')

		Expand-Archive -LiteralPath (Join-Path $tmp $archive) -DestinationPath (Join-Path $tmp 'x')
		New-Item -ItemType Directory -Force -Path $installDir | Out-Null
		$exe = Join-Path $installDir 'jaira.exe'
		# A running jaira.exe - the board, or the update check it leaves running
		# in the background - cannot be overwritten, but it can be renamed. Move
		# it aside as .old-<pid>, the name `jaira self upgrade` uses and sweeps,
		# and remove what earlier runs left there once nothing holds it.
		Get-ChildItem -LiteralPath $installDir -Filter 'jaira.exe.old-*' |
			Remove-Item -Force -ErrorAction SilentlyContinue
		if (Test-Path -LiteralPath $exe) {
			Move-Item -LiteralPath $exe -Destination "$exe.old-$PID" -Force
		}
		Copy-Item -LiteralPath (Join-Path $tmp 'x\jaira.exe') -Destination $exe
		Remove-Item -LiteralPath "$exe.old-$PID" -Force -ErrorAction SilentlyContinue
	} finally {
		Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue
	}

	Write-Host "installed $(& $exe --version) to $exe"
	if (Add-JairaUserPath $installDir) {
		Send-JairaEnvChange
		$env:Path = $env:Path.TrimEnd(';') + ';' + $installDir
		Write-Host "added $installDir to your user PATH - new terminals pick it up"
	}

	# One command writes the role skills and the cockpit, so it is one step. A
	# failure is printed as it came and does not undo the install above.
	if ($env:JAIRA_NO_ROLES -eq "1") {
		Write-Host 'skipped the skills and the Cockpit (JAIRA_NO_ROLES=1) - run: jaira roles install --global'
		return
	}
	Write-Host 'Setting up the skills and the Cockpit...'
	# 'Stop' would turn the stderr of a native command into a terminating
	# error in Windows PowerShell 5.1 and skip the failure branch below.
	$ErrorActionPreference = 'Continue'
	$out = & $exe roles install --global 2>&1 | Out-String
	$code = $LASTEXITCODE
	$ErrorActionPreference = 'Stop'
	if ($code -eq 0) {
		Write-Host ('OK ' + (($out -split "`r?`n" | Where-Object { $_ -match ' written, ' } | Select-Object -Last 1)))
		Write-Host 'active now: the jaira role skills and the Cockpit in ~\.claude\skills'
		Write-Host 'the Cockpit loads in your next Claude Code session'
	} else {
		Write-Host 'FAILED to set up the skills and the Cockpit:'
		Write-Host $out
		Write-Host 'the binary is installed; retry with: jaira roles install --global'
	}
}

Install-Jaira
