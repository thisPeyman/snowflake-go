# Snowflake ID Generator

A simple and efficient Snowflake ID generator implemented in Go. This generator creates unique, time-ordered IDs based on the Snowflake pattern, which consists of:

- **Timestamp**: Millisecond precision timestamp (relative to an epoch).
- **Machine ID**: Unique identifier for each machine or instance.
- **Sequence Number**: A sequence number to avoid collisions during the same millisecond.

## Features

- Generates unique 64-bit IDs that are sortable by creation time.
- Supports millisecond precision for timestamps.
- Built with concurrency safety using `sync.Mutex`.
- Customizable machine ID for distributed systems.
- Error handling for invalid machine IDs and timestamp issues.

## How It Works

The ID is composed of the following parts:

- **41 bits for the timestamp**: Millisecond precision, making the ID unique for at least 69 years.
- **10 bits for the machine ID**: Unique per machine or instance in the distributed system (up to 1024 machines).
- **12 bits for the sequence number**: Unique within a millisecond, with support for up to 4096 IDs per millisecond.

The format of the generated ID is:

```
| 41 bits Timestamp | 10 bits Machine ID | 12 bits Sequence |
```

## Installation

To install this package, use the following command:

```bash
go get github.com/thisPeyman/snowflake
```

## Usage

### Creating a New Snowflake Generator

You can create a new Snowflake generator by providing a unique `machineID` for your instance.

```go
package main

import (
	"fmt"
	"log"
	"github.com/thisPeyman/snowflake-go"
)

func main() {
	// Replace with a unique machine ID (must be between 0 and 1023)
	generator, err := snowflake.New(1)
	if err != nil {
		log.Fatalf("Error creating Snowflake generator: %v", err)
	}

	// Generate a unique ID
	id, err := generator.GenerateID()
	if err != nil {
		log.Fatalf("Error generating ID: %v", err)
	}

	fmt.Printf("Generated ID: %d\n", id)
}
```

### Error Handling

- **ErrInvalidMachineID**: This error is returned when the provided `machineID` is greater than the allowed limit.
- **ErrTimestampIsInvalid**: This error is returned if the timestamp is earlier than the last generated ID’s timestamp, indicating a system clock issue.

### Adjusting the Epoch

The default epoch is set to `1740304858525`. You can change it by modifying the `epoch` constant in the source code.

## Example Output

When you run the code, you'll get a unique ID that is sortable based on time:

```
Generated ID: 808586712435456000
```

## Concurrency

This implementation is thread-safe. Multiple goroutines can generate IDs concurrently without issues, thanks to the use of `sync.Mutex`.

## Performance Considerations

While this implementation uses `sync.Mutex` to ensure thread safety, under heavy load, acquiring locks could become a performance bottleneck. For high-performance scenarios, consider alternatives such as `sync/atomic` or other lock-free techniques for managing the sequence and timestamp.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.