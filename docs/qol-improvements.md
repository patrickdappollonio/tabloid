## Cleaning up extra whitespace

By default, `tabloid` will remove extra whitespace from the original output. The goal here is to provide human-readable outputs and, as seen above, `grep` or `awk` might work, but the additional whitespaces between columns are kept from the original. `tabloid` will reorganize the columns to maintain the 3-space padding between columns based on its data.
