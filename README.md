# SumHasher

A simple hasher which sums all bytes in a file.

See [CLI tool](cmd) as an example.

## Algorithm:

* for each read byte:
  * sum += (read byte + 1)
  * bytes_read++
* return max(uint64) - (sum + bytes_read)
