# Tickr — Multi-Agent RL Market Simulation

**An from-scratch financial market simulator exploring market microstructure, algorithmic trading, and deep reinforcement learning.**

Tickr is a research and portfolio project to build a financial market simulator from the ground up, starting with a Go-based limit order book and progressing toward a complete matching engine, simulated market participants, and reinforcement learning trading agents.

The long-term goal is to investigate how trading agents learn in a controlled market environment, how their strategies compare against simple baselines, and how market dynamics change when different types of participants interact.

**The guiding principle:** Build the market mechanics yourself, implement the learning algorithms, and evaluate results honestly.

> **Project status:** Early development. The foundational limit order book is implemented and tested. Matching, simulation, and reinforcement learning components are planned.

---

## Project Goals

Tickr is being developed in three layers.

### 1. Limit Order Book & Matching Engine — Go

The market infrastructure begins with a limit order book that maintains separate collections of buy and sell orders.

**Implemented:**

* Buy and sell order representations.
* Integer-based prices denominated in paise.
* Bids sorted by descending price.
* Asks sorted by ascending price.
* FIFO ordering for orders at equal prices.
* Validation of order side, price, and quantity.
* Error handling for invalid orders.
* Defensive copies returned by order-book getters.
* Automated tests for ordering, validation, and data isolation.

**Planned:**

* Order matching using price-time priority.
* Partial and complete order fills.
* Order cancellation.
* Trade generation and execution records.
* Edge-case testing and performance benchmarks.

### 2. Market Simulation

Once the matching engine is implemented, simulated market participants will generate orders and interact through the exchange.

Planned participants include:

* **Noise traders:** Submit orders with stochastic behaviour.
* **Momentum traders:** Trade based on recent price movements.
* **Market makers:** Quote bid and ask prices while managing inventory exposure.

The simulator will provide a controlled environment for studying order flow, price formation, liquidity, and trading-agent behaviour.

### 3. Reinforcement Learning — Python & PyTorch

The long-term objective is to train an RL agent to make trading decisions inside the simulated market.

Planned components:

* A Gymnasium-compatible trading environment.
* A Deep Q-Network (DQN) implemented using PyTorch.
* Random-policy and rule-based baselines.
* A Proximal Policy Optimization (PPO) implementation.
* An experiment comparing feed-forward neural networks with LSTM-based state representations.
* A multi-agent extension exploring interactions among different trading strategies.

The RL agent will interact with the simulator through a defined interface, keeping the market mechanics independent of the trading policy.

---

## Architecture

The following diagram describes the **target architecture**, not the current implementation.

```text
                 Go Market Simulator
        ┌─────────────────────────────────┐
        │                                 │
        │  Limit Order Book               │
        │          │                      │
        │          ▼                      │
        │  Matching Engine                │
        │          ▲                      │
        │          │                      │
        │  Simulated Market Participants  │
        │  (Noise, Momentum, Market Maker)│
        │                                 │
        └────────────────┬────────────────┘
                         │
                  Market State / Trades
                         │
                         ▼
                Python RL Environment
        ┌─────────────────────────────────┐
        │                                 │
        │  Gymnasium Environment          │
        │          │                      │
        │          ▼                      │
        │  RL Agent (DQN / PPO)           │
        │          │                      │
        │          ▼                      │
        │  Actions / Orders               │
        │                                 │
        └────────────────┬────────────────┘
                         │
                         ▼
                Evaluation & Analysis
        ┌─────────────────────────────────┐
        │  Baseline Comparisons           │
        │  Training Metrics               │
        │  Experiment Tracking            │
        │  Performance Visualisation      │
        └─────────────────────────────────┘
```

The intended design separates market mechanics from trading policy. This will allow different agents to interact with the same market without rewriting the underlying matching engine.

---

## Tech Stack

| Component           | Technology        | Status      |
| ------------------- | ----------------- | ----------- |
| Order book          | Go                | Implemented |
| Matching engine     | Go                | Planned     |
| Market simulation   | Go                | Planned     |
| RL environment      | Python, Gymnasium | Planned     |
| RL algorithms       | Python, PyTorch   | Planned     |
| Agent evaluation    | Python            | Planned     |
| Experiment tracking | MLflow            | Planned     |
| Results dashboard   | Streamlit         | Planned     |
| Containerisation    | Docker            | Planned     |

The initial implementation prioritises correctness and simplicity. Additional infrastructure will be introduced when the project requires it.

---

## Project Status

| Phase       | Focus                                      | Status       |
| ----------- | ------------------------------------------ | ------------ |
| Week 1      | Limit order book and validation            | Complete     |
| Week 2      | Matching engine                            | Upcoming     |
| Week 3      | Market simulation and trading bots         | Upcoming     |
| Week 4      | Gymnasium environment and baselines        | Upcoming     |
| Week 5      | DQN implementation                         | Upcoming     |
| Week 6      | Initial training and evaluation milestone  | Upcoming     |
| Weeks 7–8   | PPO and improved evaluation                | Upcoming     |
| Weeks 9–10  | MLP vs. LSTM experiment                    | Upcoming     |
| Weeks 11–12 | Robustness, analysis, and final evaluation | Upcoming     |
| Extension   | Multi-agent market simulation              | Stretch goal |

This roadmap is a development plan, not a claim that these components have already been completed.

---

## Getting Started

### Prerequisites

* Go toolchain compatible with the project's `go.mod`.
* Git.

Python and PyTorch dependencies will be documented when the RL components are introduced.

### Clone the Repository


```bash
git clone https://github.com/prachie05/tickr-multi-agent-rl-market-simulation.git
cd tickr-multi-agent-rl-market-simulation
```

### Run the Order Book Tests

```bash
cd engine
go test -v ./...
```

### Run Static Analysis

```bash
go vet ./...
```

### Format the Go Code

```bash
gofmt -w *.go
```

**Current scope:** These commands run the implemented Go order-book component. Instructions for launching the matching engine, simulation, and RL training pipeline will be added as those components become available.

---

## Repository Structure

The intended structure is:

```text
tickr/
├── engine/          # Go order book and matching engine
├── simulation/      # Market simulation and participant strategies
├── rl/              # Python RL algorithms and environment
├── experiments/     # Experiment configurations and results
├── dashboard/       # Results visualisation
└── docs/            # Architecture and technical notes
```

Directories and components will be added as development progresses.

---

## Evaluation Philosophy

A trading agent should not be considered successful simply because its training reward increases.

Tickr aims to evaluate agents using reproducible experiments and meaningful comparisons.

* **Baselines first:** Compare learned agents against random and rule-based strategies.
* **Separate evaluation scenarios:** Evaluate agents on market scenarios that were not used during training.
* **Multiple random seeds:** Report results across repeated runs using summary statistics such as mean and standard deviation.
* **Risk-aware evaluation:** Consider transaction costs, drawdowns, inventory exposure, and other relevant risk measures alongside returns.
* **Honest reporting:** Report negative or inconclusive results rather than presenting only successful experiments.
* **Reproducibility:** Record experiment configurations and random seeds so results can be independently checked.

These are evaluation principles for the project; the complete evaluation pipeline is still under development.

---

## Roadmap

* [x] Implement the foundational limit order book.
* [x] Maintain price priority and FIFO ordering.
* [x] Validate orders and reject invalid inputs.
* [x] Add automated tests for core order-book behaviour.
* [ ] Implement matching, partial fills, and cancellation.
* [ ] Build the event-driven market simulation.
* [ ] Implement simulated market participants.
* [ ] Create the Gymnasium trading environment.
* [ ] Implement and evaluate a DQN agent.
* [ ] Add a PPO baseline.
* [ ] Compare feed-forward and LSTM-based state representations.
* [ ] Evaluate robustness across different market conditions.
* [ ] Explore multi-agent trading dynamics.

---

## Disclaimer

Tickr is an educational research and software engineering project. It does not connect to live financial markets, execute real trades, or provide financial advice.

All markets, participants, and trading outcomes generated by the simulator are synthetic. Results obtained in simulation should not be interpreted as evidence of profitability in real markets.

---

