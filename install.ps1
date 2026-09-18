[CmdletBinding()]
param(
    [string]$ConfigDir = '',
    [string]$Harness = 'claude',
    [switch]$DryRun
)
$ErrorActionPreference = 'Stop'
$packageRoot = $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($packageRoot)) { throw 'Extraia o pacote e execute este arquivo com -File.' }
$packageBinary = Join-Path $packageRoot 'bin\token-saver.exe'
$checksumFile = Join-Path $packageRoot 'SHA256SUMS.txt'
if (-not (Test-Path -LiteralPath $packageBinary -PathType Leaf)) { throw 'Pacote incompleto: executavel nao encontrado. Extraia o ZIP inteiro.' }
$checksumLines = @(Get-Content -LiteralPath $checksumFile | Where-Object { $_ -match '^[a-fA-F0-9]{64}  bin/token-saver\.exe$' })
if ($checksumLines.Count -ne 1) { throw 'Manifesto de integridade invalido.' }
$expectedHash = ($checksumLines[0] -split '  ')[0]
$actualHash = (Get-FileHash -LiteralPath $packageBinary -Algorithm SHA256).Hash
if ($actualHash -ne $expectedHash) { throw 'O pacote esta corrompido. Baixe novamente a versao publicada.' }
$binaryArgs = @('install', '--harness', $Harness)
if ($ConfigDir) { $binaryArgs += @('--config-dir', $ConfigDir) }
if ($DryRun) { $binaryArgs += '--dry-run' }
Write-Host "Instalando Token Saver para o harness: $Harness"
& $packageBinary @binaryArgs
if ($LASTEXITCODE -ne 0) { throw "A instalacao nao foi concluida (codigo $LASTEXITCODE)." }
