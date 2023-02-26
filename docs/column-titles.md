# Column Handling

- [Column Handling](#column-handling)
  - [Column titles always on by default](#column-titles-always-on-by-default)
  - [Column title normalization](#column-title-normalization)
  - [Column selection and reordering](#column-selection-and-reordering)
  - [Limitations](#limitations)

## Column titles always on by default

The column titles are always on by default, so you don't have to worry about manually selecting them. Want them off? Use `--no-titles`.

## Column title normalization

In order to allow query expressions, titles are normalized: any non alphanumeric characters are removed, with the exception of `-` (dash) which is converted to underscore, and spaces are also replaced with underscores. This convention can be used both for the query expressions as well as the column selector.

In the column selector, you can also use the original column name as well in both uppercase and lowercase format.

An example conversion will be:

```diff
- NAME (PROVIDED)
+ name_provided
```

Moreover, if you prefer to see the columns before working with them, you can use `--titles-only` to print a list of titles and exit. For example, consider the following fictitional input file called `pods.txt`:

```
NAMESPACE       NAME (PROVIDED)                      READY   STATUS            RESTARTS         AGE
argocd          argocd-application-controller-0      1/1     Running           0                8d
argocd          argocd-dex-server-6dcf645b6b-nf2xb   1/1     Running           0                12d
argocd          argocd-redis-5b6967fdfc-48z9d        1/1     Running           0                14d
argocd          argocd-repo-server-7598bf5999-jfqlt  1/1     Running           0                12d
argocd          argocd-server-79f9bc9b44-5fdsp       1/1     Running           0                12d
kube-system     fluentbit-gke-qx76z                  2/2     Running           3 (2d ago)       8d
kube-system     fluentbit-gke-s2f82                  0/1     CrashLoopBackOff  592 (3m33s ago)  1h
kube-system     fluentbit-gke-wm55d                  2/2     Running           0                8d
kube-system     gke-metrics-agent-5qzdd              1/1     Running           0                200d
kube-system     gke-metrics-agent-95vkn              1/1     Running           0                200d
kube-system     gke-metrics-agent-blbbm              1/1     Running           0                8d
```

You can use `--titles-only` to print the titles and exit:

```bash
$ cat pods.txt | tabloid --titles-only
NAMESPACE
NAME (PROVIDED)
READY
STATUS
RESTARTS
AGE
```

You can also combine `--titles-only` with `--titles-normalized` to print the titles post-normalization for expressions:

```bash
$ cat pods.txt | tabloid --titles-only --titles-normalized
namespace
name_provided
ready
status
restarts
age
```

## Column selection and reordering

By default, all columns are shown exactly as shown by the original. However, if one or more columns are provided -- either via the `--column` parameter using comma-separated values, or by repeating `--column` as many times as needed -- then only those columns are shown, in the order they are received.

## Limitations

* Column names must be unique.
* Column values are always strings [unless processed by a built-in function](expressions.md#expression-functions) -- this means it's not possible to perform math comparisons yet.
* The `--expr` parameter must be quoted depending on your terminal.
* The input must adhere to Go's `tabwriter` using 2 or more spaces between columns minimum (this is true for both `docker` and `kubectl`).
* Due to the previous item, column names must not contain 2+ consecutive spaces, otherwise they are treated as multiple columns, potentially breaking parsing.
