from pathlib import Path
import subprocess
base=Path('.parley-runtime/quota-implementation/fixup-2/checks')
common=r'''package main
import("bytes";"fmt";"os";"path/filepath";"strings";"parley-deck-cli/internal/app";"parley-deck-cli/internal/protocol")
func must(e error){if e!=nil{panic(e)}}
func main(){root:=os.Args[1];must(os.MkdirAll(root,0700)); must(protocol.InitWorkspace(root)); idea:=setup(root);p:=filepath.Join(idea.Path,"00-prompt.md");must(os.WriteFile(filepath.Join(idea.Path,"consensus.md"),[]byte("---\nidea: "+idea.Slug+"\ndrafted-by: a\n---\n\n## Signoffs\n"),0600));
edit:=func(from,to string){raw,e:=os.ReadFile(p);must(e);next:=strings.Replace(string(raw),from,to,1);if next==string(raw){panic("edit missed")};must(os.WriteFile(p,[]byte(next),0600))}
check:=func(label,agent string){var out,err bytes.Buffer; code:=app.Run([]string{"consensus","signoff","--dir",root,"--agent",agent,"--status","accept",idea.Slug},&out,&err);fmt.Printf("%s exit=%d stdout=%s stderr=%s\n",label,code,out.String(),err.String());if code!=0{panic("CLI compatibility failed")}}
edit("participants: [a, b, c, d]","participants: [a, b, c]\nexcluded: [d — unavailable — quota exhausted — confirmed 2026-10-04]");check("excluded reason with em dash","a");
edit("participants: [a, b, c]","participants: [a, b, c, d]");check("known return plain participants edit","d");
raw:="---\nagent: e\nidea: "+idea.Slug+"\nround: 1\n---\n";for _,h:=range []string{"Summary","Proposed approach","Concerns / open questions","Risks","Existing alternatives"}{raw+="\n## "+h+"\nRead prior proposals and join deliberation from round 2.\n"};must(os.WriteFile(filepath.Join(idea.Path,"round-01/e.md"),[]byte(raw),0600));
edit("participants: [a, b, c, d]","participants: [a, b, c, d, e]");check("late round-1 and plain participants join","e");verify(idea.Path)
}
'''
old='''package main
import "parley-deck-cli/internal/protocol"
func setup(root string) protocol.IdeaStatus {i,e:=protocol.CreateIdeaFull(root,"cycle2 differential",[]string{"a","b","c","d"},nil,"deliberation","");must(e);return i}
func verify(dir string){}
'''
new='''package main
import("fmt";"parley-deck-cli/internal/protocol";"parley-deck-cli/internal/quota";"parley-deck-cli/internal/runmanifest")
func setup(root string) protocol.IdeaStatus {p:=quota.Policy{Enabled:false,Scope:quota.KickoffAndMidIdea};i,k,e:=protocol.CreateIdeaWithQuota(root,"cycle2 differential",[]string{"a","b","c","d"},nil,"deliberation","","run-1",&p,nil);must(e);m:=runmanifest.New(runmanifest.Options{Root:root,RunID:"run-1",IdeaSlug:i.Slug,Participants:k.Participants,QuotaKickoff:k});must(runmanifest.Write(root,"run-1",m));return i}
func verify(dir string){h,e:=quota.ReadHistory(dir);must(e);if h.Revision!=3 {panic("missing immutable manual history")};for _,b:=range h.Batches{if b.Owner==nil||b.Owner.ManualPrompt==""||b.Owner.Authority!=nil{panic("manual imported as owner authority")}};fmt.Println("three immutable manual revisions; no owner-confirmed authority")}
'''
for path,setup in [(base/'differential-current',new),(Path('/tmp/quota-cycle2-prechange-27e42b8')/'.cycle2-differential',old)]:
 path.mkdir(exist_ok=True);(path/'main.go').write_text(common);(path/'setup.go').write_text(setup)
(base/'differential-common.go.txt').write_text(common);(base/'differential-baseline-setup.go.txt').write_text(old)
# Reviewer input reproduction: only change the producer output locator.
src=Path('.parley-runtime/claude1-r4/zretry/main.go');dst=base/'reviewer-zretry';dst.mkdir(exist_ok=True)
(dst/'main.go').write_text(src.read_text().replace('.parley-runtime/claude1-r4/zretry/B-aggregated-changing-countdown.stderr',str(base/'reviewer-B-aggregated-changing-countdown.stderr')))
