# Stellar Chronodynamics
Public, standalone, science-anchored procedural worldbuilding bake-off. Unpublished story scenarios are isolated in SemperSupra/stellar-chronodynamics-private.

## Run
Requires Go 1.22+.

1. go test ./...
2. go test -race ./...
3. go vet ./...
4. go run ./cmd/bakeoff -input fixtures/three_systems.json -out .local/experiment

## Evidence status
M0 runs a deliberately SYNTHETIC toy stellar flux model to qualify temporal overlap semantics and deterministic execution. It does not establish that a real planet is habitable or inhabited.

Work plan and execution authority: https://github.com/SemperSupra/stellar-chronodynamics/issues/1

Science candidates: NASA Exoplanet Archive TAP, HZ_evolution, VPLanet, REBOUND, QuantHab. No upstream scientific model qualified yet.
