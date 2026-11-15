package main

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "flag"
 "fmt"
 "math"
 "os"
 "path/filepath"
 "sort"
)

type Point struct { T float64 `json:"epoch_gyr"`; L float64 `json:"luminosity_solar"` }
type Planet struct { ID string `json:"id"`; A float64 `json:"semi_major_axis_au"`; Min float64 `json:"flux_min"`; Max float64 `json:"flux_max"` }
type System struct { ID string `json:"id"`; Track []Point `json:"luminosity_track"`; Planets []Planet `json:"planets"` }
type Input struct { Version string `json:"schema_version"`; Synthetic bool `json:"synthetic"`; Start float64 `json:"epoch_start_gyr"`; End float64 `json:"epoch_end_gyr"`; Step float64 `json:"step_gyr"`; Systems []System `json:"systems"` }
type Window struct { System string `json:"system_id"`; Planet string `json:"planet_id"`; Start float64 `json:"start_gyr"`; End float64 `json:"end_gyr"` }
type Overlap struct { A string `json:"system_a"`; B string `json:"system_b"`; Start float64 `json:"start_gyr"`; End float64 `json:"end_gyr"` }
type Report struct { Grade string `json:"evidence_grade"`; InputSHA string `json:"input_sha256"`; Windows []Window `json:"windows"`; Overlaps []Overlap `json:"overlaps"`; Limits []string `json:"limitations"` }
func finite(v float64) bool{return !math.IsNaN(v)&&!math.IsInf(v,0)}
func validate(v Input)error{
 if v.Version!="0.1"||!v.Synthetic{return errors.New("M0 requires schema 0.1 and synthetic=true; observed records not admitted")}
 if !finite(v.Start)||!finite(v.End)||!finite(v.Step)||v.End<=v.Start||v.Step<=0||(v.End-v.Start)/v.Step>100000{return errors.New("invalid time range")}
 if len(v.Systems)==0||len(v.Systems)>100{return errors.New("invalid system count")}
 ids:=map[string]bool{}
 for _,s:=range v.Systems{
  if s.ID==""||ids[s.ID]||len(s.Track)<2||len(s.Planets)==0{return fmt.Errorf("invalid system %q",s.ID)}
  ids[s.ID]=true
  if s.Track[0].T>v.Start||s.Track[len(s.Track)-1].T<v.End{return fmt.Errorf("incomplete track %q",s.ID)}
  prev:=math.Inf(-1)
  for _,p:=range s.Track {if !finite(p.T)||!finite(p.L)||p.T<=prev||p.L<=0{return fmt.Errorf("invalid track %q",s.ID)};prev=p.T}
  seen:=map[string]bool{}
  for _,p:=range s.Planets{
   if p.ID==""||seen[p.ID]||!finite(p.A)||!finite(p.Min)||!finite(p.Max)||p.A<=0||p.Min<=0||p.Max<=p.Min{return fmt.Errorf("invalid planet %q",p.ID)}
   seen[p.ID]=true
  }
 }
 return nil
}
func lum(a []Point,t float64)float64{
 if t<=a[0].T{return a[0].L}
 for i:=1;i<len(a);i++{if t<=a[i].T{x,y:=a[i-1],a[i];return x.L+(y.L-x.L)*(t-x.T)/(y.T-x.T)}}
 return a[len(a)-1].L
}
func round(n float64)float64{return math.Round(n*1e6)/1e6}
func run(v Input,raw []byte)(Report,error){
 if err:=validate(v);err!=nil{return Report{},err}
 h:=sha256.Sum256(raw)
 o:=Report{Grade:"TOY_SYNTHETIC_ONLY",InputSHA:hex.EncodeToString(h[:]),Windows:[]Window{},Overlaps:[]Overlap{},Limits:[]string{"Synthetic luminosity and fixed flux thresholds, not scientific habitable-zone estimates","No observed star, atmosphere, orbit evolution, life or civilization inference","Interval resolution depends on step; no biological claim"}}
 n:=int(math.Ceil((v.End-v.Start)/v.Step))
 for _,s:=range v.Systems{for _,p:=range s.Planets{
  active:=false;start:=0.0
  for i:=0;i<n;i++{
   a:=v.Start+float64(i)*v.Step;b:=math.Min(v.End,a+v.Step)
   flux:=lum(s.Track,(a+b)/2)/(p.A*p.A)
   ok:=flux>=p.Min&&flux<=p.Max
   if ok&&!active{active=true;start=a}
   if !ok&&active{o.Windows=append(o.Windows,Window{s.ID,p.ID,round(start),round(a)});active=false}
  }
  if active{o.Windows=append(o.Windows,Window{s.ID,p.ID,round(start),round(v.End)})}
 }}
 sort.Slice(o.Windows,func(i,j int)bool{a,b:=o.Windows[i],o.Windows[j];if a.System!=b.System{return a.System<b.System};if a.Planet!=b.Planet{return a.Planet<b.Planet};return a.Start<b.Start})
 for i,a:=range o.Windows{for _,b:=range o.Windows[i+1:]{if a.System==b.System{continue};s:=math.Max(a.Start,b.Start);e:=math.Min(a.End,b.End);if e>s{o.Overlaps=append(o.Overlaps,Overlap{a.System,b.System,round(s),round(e)})}}}
 return o,nil
}
func execute(input,out string)error{
 raw,err:=os.ReadFile(input);if err!=nil{return err}
 var v Input
 if err=json.Unmarshal(raw,&v);err!=nil{return err}
 r,err:=run(v,raw);if err!=nil{return err}
 if err=os.MkdirAll(out,0755);err!=nil{return err}
 b,err:=json.MarshalIndent(r,"","  ");if err!=nil{return err}
 if err=os.WriteFile(filepath.Join(out,"result.json"),append(b,'\n'),0644);err!=nil{return err}
 summary:=fmt.Sprintf("# Synthetic-only baseline\n\nEvidence: TOY_SYNTHETIC_ONLY\n\nInput SHA256: %s\n\nWindows: %d\nInter-system overlaps: %d\n\nThis is NOT an exoplanet habitability model.\n",r.InputSHA,len(r.Windows),len(r.Overlaps))
 if err=os.WriteFile(filepath.Join(out,"summary.md"),[]byte(summary),0644);err!=nil{return err}
 fmt.Printf("PASS synthetic baseline windows=%d overlaps=%d\n",len(r.Windows),len(r.Overlaps))
 return nil
}
func main(){
 in:=flag.String("input","fixtures/three_systems.json","synthetic fixture");out:=flag.String("out",".local/experiment","evidence directory");flag.Parse()
 if err:=execute(*in,*out);err!=nil{fmt.Fprintln(os.Stderr,"FAIL:",err);os.Exit(1)}
}
