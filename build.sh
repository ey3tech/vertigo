echo "formatting..."
gofmt -s -w *.go && echo "formatting complete..." || echo "formatting failed"
if go get github.com/ey3tech/vertigo; then
    echo "successfully installed vertigo"
    
else
    echo "failed to install vertigo"
fi