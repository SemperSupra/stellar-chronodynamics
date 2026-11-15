package main
import ("encoding/json";"reflect";"testing")
func sample()Input{return Input{Version:"0.1",Synthetic:true,Start:0,End:4,Step:0.1,Systems:[]System{
{ID:"A",Track:[]Point{{T:0,L:.5},{T:4,L:1.5}},Planets:[]Planet{{ID:"p1",A:1,Min:.35,Max:1.1}}},
{ID:"B",Track:[]Point{{T:0,L:.4},{T:4,L:.8}},Planets:[]Planet{{ID:"p2",A:1,Min:.35,Max:1.1}}},
}}}
func TestDeterministicOverlap(t *testing.T){v:=sample();raw,_:=json.Marshal(v);a,e:=run(v,raw);if e!=nil{t.Fatal(e)};b,e:=run(v,raw);if e!=nil{t.Fatal(e)};if !reflect.DeepEqual(a,b)||len(a.Overlaps)==0||a.Grade!="TOY_SYNTHETIC_ONLY"{t.Fatal("incorrect baseline")}}
func TestRejectObserved(t *testing.T){v:=sample();v.Synthetic=false;if validate(v)==nil{t.Fatal("non synthetic admitted")}}
func TestRejectBadPhysics(t *testing.T){v:=sample();v.Systems[0].Planets[0].A=0;if validate(v)==nil{t.Fatal("zero radius")};v=sample();v.Systems[0].Track[1].T=0;if validate(v)==nil{t.Fatal("time order")}}
