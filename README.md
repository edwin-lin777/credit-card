# Canadian Credit Card Rewards Simulator

A concurrent Monte Carlo simulator written in Go that compares Canadian credit-card wallets using category-specific cashback rates, spending caps, and annual fees.

The project focuses on backend engineering: a configurable rules engine, deterministic synthetic data generation, a goroutine worker pool, CPU profiling, and reproducible Docker/CI workflows.

## Features

- **Cap-aware rewards:** base and bonus rates, monthly/annual spending caps, shared category caps, spillover to lower rates, and annual fees.
- **Wallet comparison:** evaluates all unique two- and three-card combinations from a YAML card catalog.
- **Synthetic households:** generates monthly spending across eight categories using selected 2023 Statistics Canada household spending figures as anchors.
- **Concurrent Monte Carlo:** distributes household evaluations across a bounded goroutine worker pool, with worker-local totals to avoid contention.
- **Precomputed wallet rules:** sorts reward offers once per wallet, rather than once per household.
- **Reproducible results:** configurable seed, command-line flags, ranked JSON output, multistage Docker build, and GitHub Actions checks.

## Example output

With 12 modeled cards, the simulator evaluates **286 wallets** (66 two-card combinations and 220 three-card combinations).

An example local run, after the latest category-rule corrections, used 1,000 synthetic households, four workers, and seed 42:

~~~text
Cards: 12
Wallet combinations: 286
Households: 1000
Workers: 4
Elapsed time: 795.4225ms

Top 10 wallets:
1. $675.88 | Rogers Red World Elite Mastercard + Scotia Momentum Visa Infinite Plus Card + Simplii Financial Cash Back Visa
2. $668.45 | Capital One Quicksilver World Elite Mastercard + Rogers Red World Elite Mastercard + Simplii Financial Cash Back Visa
3. $668.45 | Capital One Quicksilver World Mastercard + Rogers Red World Elite Mastercard + Simplii Financial Cash Back Visa
...
Results saved to results.json
~~~

Rewards shown are **average simulated annual net cashback after modeled annual fees**, not personalized financial recommendations. Results depend on the spending assumptions and card configurations.

## Quick start

**Requirements:** Go 1.27.1 or later. Run from the repository root so the program can load `cards/*.yaml`.

~~~bash
git clone https://github.com/edwin-lin777/credit-card.git
cd credit-card
go run ./cmd/study -households 1000 -workers 4 -seed 42 -output results.json
~~~

For a 50,000-household study:

~~~bash
go run ./cmd/study -households 50000 -workers 16 -seed 42 -output results-50k.json
~~~

| Flag | Default | Description |
| --- | --- | --- |
| `-households` | 1000 | Number of synthetic households |
| `-workers` | Go's `GOMAXPROCS` | Number of concurrent workers |
| `-seed` | 42 | Seed for repeatable household generation |
| `-output` | `results.json` | Destination for ranked JSON results |

The JSON report includes the run configuration, elapsed time, number of cards and wallets, and the ranked card combinations with average net rewards in integer cents.

### Docker

Build the image:

~~~bash
docker build -t credit-card .
~~~

Run it and save the JSON output to the current directory.

PowerShell:

~~~powershell
docker run --rm -v "${PWD}:/output" credit-card -households 1000 -workers 4 -seed 42 -output /output/results.json
~~~

Bash:

~~~bash
docker run --rm -v "$(pwd):/output" credit-card -households 1000 -workers 4 -seed 42 -output /output/results.json
~~~

The bind mount is important: without it, the results file remains inside the temporary container and is removed with `--rm`.

## Architecture

~~~text
cards/*.yaml ───► Catalog loader ───► Generate 2-/3-card wallets
                                             │
                                             ▼
                                        PrepareWallet
                                      (sort offers once)
                                             │
Statistics Canada anchors ───► Generate household profiles
                               (seed + household ID)
                                             │
                                             ▼
                                      Goroutine worker pool
                                      ├─ Worker-local totals
                                      ├─ Cap-aware allocation
                                      └─ Evaluate each wallet
                                             │
                                             ▼
                                      Aggregate + rank
                                             │
                                             ▼
                                      Console + JSON
~~~

### Reward rules

Credit-card rules are stored as YAML, not hard-coded into the evaluation engine. A group can share a monthly or annual spending cap among several reward categories.

~~~yaml
id: example-card
name: Example Card
annual_fee_cents: 12000
base_rate: 0.01

groups:
  - rates:
      groceries: 0.04
      bills: 0.04
    cap:
      amount_cents: 2500000
      period: annual
~~~

In this example, `2500000` cents represents **$25,000 of spending eligible for bonus rewards** per year, shared between the two categories; it is not a $25,000 cashback payout.

### Wallet allocation

`PrepareWallet` precomputes each wallet's cashback opportunities, sorted by rate. `EvaluatePreparedWallet` processes monthly spending in that order, tracking monthly and annual cap usage and assigning eligible spending to cards. The engine then calculates earned cashback and subtracts annual fees.

The strategy is **greedy**, not a proof of globally optimal reward allocation. Annual-cap decisions can depend on spending later in the year.

### Synthetic households

`internal/households` uses selected [Statistics Canada 2023 Survey of Household Spending](https://www150.statcan.gc.ca/n1/daily-quotidien/250521/dq250521a-eng.htm) figures as category anchors. It models groceries, dining, gas, transit, bills, retail, travel, and other purchases.

For each household, annual category spending is sampled using an arbitrary uniform multiplier between **0.6 and 1.4** times the corresponding anchor, then distributed over 12 months using randomized weights. A fixed seed and household ID generate repeatable spending profiles.

Money is represented in **integer cents**, including card fees, spending amounts, and accumulated rewards.

### Concurrent study

`internal/study/runner.go` sends household IDs through a buffered channel to a configurable number of worker goroutines. Each worker generates a household profile and evaluates all prepared wallets, accumulating results in a private array. After a `sync.WaitGroup` completes, the arrays are combined and wallets are ranked by total net rewards.

This avoids concurrent writes to shared totals. A test also checks that one worker and four workers produce identical reward totals for a fixed seed.

## Performance

The included catalog has 12 cards and 286 two-/three-card wallets, so a 50,000-household run performs approximately **14.3 million household-wallet evaluations**.

Local development measurements:

| Measurement | Observed value | Notes |
| --- | --- | --- |
| 1,000-household run | ~0.80 s | 4 workers; updated card categories |
| 50,000-household run | ~34.3 s | 16 workers; earlier card-rule snapshot |
| Wallet evaluation microbenchmark | ~8.3 → ~4.2 µs/op | Before/after precomputing sorted offers |

These are illustrative observations from different development stages and environments, **not a controlled head-to-head benchmark**. Rerun benchmarks on the same machine and code revision for precise comparisons.

Run the current wallet benchmark:

~~~bash
go test ./internal/rewards -run '^$' -bench BenchmarkEvaluateWallet -benchmem
~~~

Profile CPU usage:

~~~bash
go test ./internal/rewards -run '^$' -bench BenchmarkEvaluateWallet -cpuprofile cpu.prof
go tool pprof -top cpu.prof
~~~

## Tests and CI

Run all unit tests:

~~~bash
go test ./...
~~~

Run the race detector on Linux or a machine with a working CGO C compiler:

~~~bash
CGO_ENABLED=1 go test -race ./...
~~~

The test suite covers cashback rates, cap limits, spillover, shared caps, annual fees, YAML parsing, wallet counts, and worker-count determinism.

[GitHub Actions](.github/workflows/ci.yml) runs unit tests, Go race detection, and a Docker build on pushes and pull requests.

## Repository structure

~~~text
cards/                         Credit-card YAML configurations
cmd/study/                     CLI, wallet enumeration, JSON export
internal/catalog/              Card types and YAML parsing
internal/households/           Spending anchors and synthetic generation
internal/rewards/              Reward and wallet evaluation, benchmarks
internal/spending/             Spending categories and monthly profile
internal/study/                Concurrent Monte Carlo runner
.github/workflows/ci.yml       Automated tests and Docker build
Dockerfile                     Multistage container build
~~~

## Limitations

This is an **engineering simulation**, not a live credit-card comparison or financial recommendation service.

- **Spending distributions are synthetic:** the 60–140% multipliers and monthly weights are not fitted to a representative household survey. The model does not capture correlations, income, geography, or family composition.
- **Spending anchors are partial proxies:** several categories combine selected survey line items and do not measure all credit-card-eligible Canadian household spending.
- **Merchant eligibility is simplified:** the eight broad spending categories cannot model every issuer-specific merchant code, exclusion, recurring-payment requirement, or product-specific bonus category.
- **Card terms are simplified:** promotional offers, signup bonuses, redemption boosts, eligibility rules, and some customer-specific benefits are excluded. Rates and fees can change; each YAML file includes a `source_url` and `verified_on` field.
- **Optimization is heuristic:** the highest-current-rate-first allocation is not guaranteed to maximize rewards globally over the entire year.
- **Wallet scope is limited:** only two- and three-card combinations are evaluated, without credit approval or consumer suitability modeling.

The project is intended to demonstrate maintainable backend design, concurrency, reproducible simulations, and profiling—not to predict a particular household's rewards.
