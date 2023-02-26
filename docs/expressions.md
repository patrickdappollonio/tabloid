# Expressions

- [Expressions](#expressions)
  - [Powerful expression evaluator](#powerful-expression-evaluator)
  - [Expression functions](#expression-functions)
    - [`isready`, `isnotready`](#isready-isnotready)
    - [`hasrestarts`, `hasnorestarts`](#hasrestarts-hasnorestarts)
    - [`olderthan`, `olderthaneq`, `newerthan`, `newerthaneq`, `eqduration`](#olderthan-olderthaneq-newerthan-newerthaneq-eqduration)

## Powerful expression evaluator

The `--expr` parameter allows you to specify any boolean expression. `tabloid` uses [`govaluate`](https://github.com/Knetic/govaluate) for its expression evaluator and multiple options are supported, such as:

* Grouping with parenthesis
* `&&` and `||` operators
* `!=`, `==`, `>`, `<`, `>=`, `<=` operators
* And regexp-based operators such as `=~` and `!~`, based on Go's own `regexp` package

The only requirement, evaluated after parsing your expression, is that the expression must evaluate to a boolean output.

Mathematical operators do not work due to how the table is parsed: all values are strings.

## Expression functions

The following functions are available. Their parameters are the column names and potential additional values you want to pass to them. See their examples for more details.

**Note:** Expressions always use the normalized column name as an input parameter or matching value. For example, if you have a column named `NAME (PROVIDED)`, then you would use `name_provided` as the parameter name.

The following file is used for the examples below:

```bash
$ cat pods.txt
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

### `isready`, `isnotready`

Returns true for any value that matches the format `<current>/<total>`. If `<current>` is equal to `<total>`, then `isready` returns true, otherwise `isnotready` returns true.

These functions will work with columns that contains values such as `1/1` or `0/1`. A row is considered "not ready" if the `<current>` value is not equal to the `<total>` value.

**Examples:**

```bash
# Print all pods that have an amount of pods matching the expected amount
$ cat pods.txt | tabloid --expr 'isready(ready)'
NAMESPACE     NAME (PROVIDED)                       READY   STATUS    RESTARTS     AGE
argocd        argocd-application-controller-0       1/1     Running   0            8d
argocd        argocd-dex-server-6dcf645b6b-nf2xb    1/1     Running   0            12d
argocd        argocd-redis-5b6967fdfc-48z9d         1/1     Running   0            14d
argocd        argocd-repo-server-7598bf5999-jfqlt   1/1     Running   0            12d
argocd        argocd-server-79f9bc9b44-5fdsp        1/1     Running   0            12d
kube-system   fluentbit-gke-qx76z                   2/2     Running   3 (2d ago)   8d
kube-system   fluentbit-gke-wm55d                   2/2     Running   0            8d
kube-system   gke-metrics-agent-5qzdd               1/1     Running   0            200d
kube-system   gke-metrics-agent-95vkn               1/1     Running   0            200d
kube-system   gke-metrics-agent-blbbm               1/1     Running   0            8d
```

```bash
# Print all pods that have an amount of pods NOT matching the expected amount
$ cat pods.txt | tabloid --expr 'isnotready(ready)'
NAMESPACE     NAME (PROVIDED)       READY   STATUS             RESTARTS          AGE
kube-system   fluentbit-gke-s2f82   0/1     CrashLoopBackOff   592 (3m33s ago)   1h
```

### `hasrestarts`, `hasnorestarts`

Returns true for any value that matches the format `<number>` where `<number>` is a positive integer. Optionally, it also supports values whith the format `<number> (<duration> ago)`, where `<number>` is a positive integer, and `<duration>` is a Go-parseable `time.Duration` (with additional support up to days, like `kubectl`).

If `<number>` is greater than 0, then `hasrestarts` returns true, otherwise `hasnorestarts` returns true.

Column values could have formats like `5` or `5 (5d ago)`. The value within the parenthesis is ignored.

**Examples:**

```bash
# Print all pods that have had at least one restart
$ cat pods.txt | tabloid --expr 'hasrestarts(restarts)'
NAMESPACE     NAME (PROVIDED)       READY   STATUS             RESTARTS          AGE
kube-system   fluentbit-gke-qx76z   2/2     Running            3 (2d ago)        8d
kube-system   fluentbit-gke-s2f82   0/1     CrashLoopBackOff   592 (3m33s ago)   1h
```

```bash
# Print all pods that have had no restarts
$ cat pods.txt | tabloid --expr 'hasnorestarts(restarts)'
NAMESPACE     NAME (PROVIDED)                       READY   STATUS    RESTARTS   AGE
argocd        argocd-application-controller-0       1/1     Running   0          8d
argocd        argocd-dex-server-6dcf645b6b-nf2xb    1/1     Running   0          12d
argocd        argocd-redis-5b6967fdfc-48z9d         1/1     Running   0          14d
argocd        argocd-repo-server-7598bf5999-jfqlt   1/1     Running   0          12d
argocd        argocd-server-79f9bc9b44-5fdsp        1/1     Running   0          12d
kube-system   fluentbit-gke-wm55d                   2/2     Running   0          8d
kube-system   gke-metrics-agent-5qzdd               1/1     Running   0          200d
kube-system   gke-metrics-agent-95vkn               1/1     Running   0          200d
kube-system   gke-metrics-agent-blbbm               1/1     Running   0          8d
```

### `olderthan`, `olderthaneq`, `newerthan`, `newerthaneq`, `eqduration`

Utility functions to manage durations, as seen in the `kubectl` output. These functions are useful to compare durations, such as the age of a pod. You can use it to formulate queries like "all pods older or equal than 1 day".

These functions will work with columns that contains values parseable by `time.ParseDuration` (with additional support up to days, like `kubectl`).

**Examples:**

```bash
# Print all pods that are older than 8 days
$ cat pods.txt | tabloid --expr 'olderthan(age, "8d")'
NAMESPACE     NAME (PROVIDED)                       READY   STATUS    RESTARTS   AGE
argocd        argocd-dex-server-6dcf645b6b-nf2xb    1/1     Running   0          12d
argocd        argocd-redis-5b6967fdfc-48z9d         1/1     Running   0          14d
argocd        argocd-repo-server-7598bf5999-jfqlt   1/1     Running   0          12d
argocd        argocd-server-79f9bc9b44-5fdsp        1/1     Running   0          12d
kube-system   gke-metrics-agent-5qzdd               1/1     Running   0          200d
kube-system   gke-metrics-agent-95vkn               1/1     Running   0          200d
```

```bash
# Print all pods that are older or equal than 8 days
$ cat pods.txt | tabloid --expr 'olderthaneq(age, "8d")'
NAMESPACE     NAME (PROVIDED)                       READY   STATUS    RESTARTS     AGE
argocd        argocd-application-controller-0       1/1     Running   0            8d
argocd        argocd-dex-server-6dcf645b6b-nf2xb    1/1     Running   0            12d
argocd        argocd-redis-5b6967fdfc-48z9d         1/1     Running   0            14d
argocd        argocd-repo-server-7598bf5999-jfqlt   1/1     Running   0            12d
argocd        argocd-server-79f9bc9b44-5fdsp        1/1     Running   0            12d
kube-system   fluentbit-gke-qx76z                   2/2     Running   3 (2d ago)   8d
kube-system   fluentbit-gke-wm55d                   2/2     Running   0            8d
kube-system   gke-metrics-agent-5qzdd               1/1     Running   0            200d
kube-system   gke-metrics-agent-95vkn               1/1     Running   0            200d
kube-system   gke-metrics-agent-blbbm               1/1     Running   0            8d
```

```bash
# Print all pods that are newer than 8 days
$ cat pods.txt | tabloid --expr 'newerthan(age, "8d")'
NAMESPACE     NAME (PROVIDED)       READY   STATUS             RESTARTS          AGE
kube-system   fluentbit-gke-s2f82   0/1     CrashLoopBackOff   592 (3m33s ago)   1h
```

```bash
# Print all pods that are newer or equal than 8 days
$ cat pods.txt | tabloid --expr 'newerthaneq(age, "8d")'
NAMESPACE     NAME (PROVIDED)                   READY   STATUS             RESTARTS          AGE
argocd        argocd-application-controller-0   1/1     Running            0                 8d
kube-system   fluentbit-gke-qx76z               2/2     Running            3 (2d ago)        8d
kube-system   fluentbit-gke-s2f82               0/1     CrashLoopBackOff   592 (3m33s ago)   1h
kube-system   fluentbit-gke-wm55d               2/2     Running            0                 8d
kube-system   gke-metrics-agent-blbbm           1/1     Running            0                 8d
```

```bash
# Print all pods that are exactly 8 days old
$ cat pods.txt | tabloid --expr 'eqduration(age, "8d")'
NAMESPACE     NAME (PROVIDED)                   READY   STATUS    RESTARTS     AGE
argocd        argocd-application-controller-0   1/1     Running   0            8d
kube-system   fluentbit-gke-qx76z               2/2     Running   3 (2d ago)   8d
kube-system   fluentbit-gke-wm55d               2/2     Running   0            8d
kube-system   gke-metrics-agent-blbbm           1/1     Running   0            8d
```
