$file = "./platform.exe"
&go build -o $file
if ($LASTEXITCODE -eq 0) {
&$file
}
