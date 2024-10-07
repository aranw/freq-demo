# Temporal Frequency Demo

This is a Temporal.io demo that processes generated frequency data simulating the UK National Grid's 50Hz grid.

## Overview

The Frequency Demo project showcases the capabilities of Temporal.io in handling batched data processing. It simulates the frequency fluctuations of the UK National Grid, which operates at a nominal frequency of 50Hz. The project generates frequency readings, batches them, and then processes these batches using Temporal.io for workflow orchestration.

## Prerequisites

* Go 1.23 or later
* Temporal.io server

## Usage

1. Start the Temporal worker for batch processing:
   ```
   go run worker
   ```

2. Run the frequency data generator and batcher:
   ```
   go run generator/cmd/generator
   ```

## Project Structure

- `generator/`: Contains the frequency data generator and batcher
- `worker/`: Implements the Temporal.io worker for batch processing
- `workflow/`: Defines the Temporal.io workflows for orchestration

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
