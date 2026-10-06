package main
import("os";"path/filepath";"parley-deck-cli/internal/protocol")
func main(){root:=os.Args[1];if e:=protocol.InitWorkspace(root);e!=nil{panic(e)};p:=filepath.Join(root,"parley-deck/meta/version.json");if e:=os.WriteFile(p,[]byte(`{"protocolRole":"source"}`),0600);e!=nil{panic(e)}}
