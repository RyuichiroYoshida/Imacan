$ErrorActionPreference = "Stop"

$source = "../../api/openapi/project-imacan.openapi.yaml"
$temp = ".codegen.openapi.yaml"

try {
    Get-Content -Encoding UTF8 $source |
        Where-Object { $_ -notmatch '^\s*(description|summary):' } |
        Set-Content -Encoding UTF8 $temp

    go tool oapi-codegen -config oapi-codegen.yaml $temp
}
finally {
    if (Test-Path $temp) {
        Remove-Item $temp
    }
}
