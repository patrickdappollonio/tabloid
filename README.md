# `tabloid` -- your tabulated data's best friend

[![Downloads](https://img.shields.io/github/downloads/patrickdappollonio/tabloid/total?color=blue&logo=github&style=flat-square)](https://github.com/patrickdappollonio/tabloid/releases)

`tabloid` is a weekend project. The goal is to be able to **parse inputs from several command line applications like `kubectl` and `docker` that use a `tabwriter` to format their output**: this is, they write column-based outputs where the first line is the column title -- often uppercased -- and the values come below, and they're often perfectly aligned.

Here's an example: more often than not though, you want one field from that output instead of a tens-of-lines long output. So your first attempt is to resort to `grep`:

```bash
$ kubectl get pods --all-namespaces | grep frontend
team-a-apps     frontend-5c6c94684f-5kzbk                                       1/1     Running   0          8d
team-a-apps     frontend-5c6c94684f-k2d7d                                       1/1     Running   0          8d
team-a-apps     frontend-5c6c94684f-ppgkx                                       1/1     Running   0          8d
```

You have a couple of issues here:

* The first column disappeared, which holds the titles. I'm often forgetful and won't remember what each column is supposed to be. Maybe for some outputs, but not all (looking at you, `kubectl api-resources`!)
* There's some awkward space between the columns now, since the columns keep the original formatting.

We could fix the first issue by using `awk` instead of `grep`:

```bash
$ kubectl get pods --all-namespaces | awk 'NR == 1 || /frontend/'
NAMESPACE       NAME                                                            READY   STATUS    RESTARTS   AGE
team-a-apps     frontend-5c6c94684f-5kzbk                                       1/1     Running   0          8d
team-a-apps     frontend-5c6c94684f-k2d7d                                       1/1     Running   0          8d
team-a-apps     frontend-5c6c94684f-ppgkx                                       1/1     Running   0          8d
```

Much better! Now if this works for you, you can stop reading here. Chances are, you won't need `tabloid`. But if you want:

* Some more human-readable filters than `awk`
* The ability to customize the columns' order
* The ability to filter with `AND` and `OR` rules
* Or filter using regular expressions

Then `tabloid` is the right tool for you. Here's an example:

```bash
# show only pods whose name starts with `frontend` or `redis`
$ kubectl get pods --all-namespaces | tabloid --expr 'name =~ "^frontend" || name =~ "^redis"'
NAMESPACE     NAME                             READY   STATUS    RESTARTS   AGE
team-a-apps   frontend-5c6c94684f-5kzbk        1/1     Running   0          8d
team-a-apps   frontend-5c6c94684f-k2d7d        1/1     Running   0          8d
team-a-apps   frontend-5c6c94684f-ppgkx        1/1     Running   0          8d
team-a-apps   redis-follower-dddfbdcc9-9xd8l   1/1     Running   0          8d
team-a-apps   redis-follower-dddfbdcc9-l9ngl   1/1     Running   0          8d
team-a-apps   redis-leader-fb76b4755-6t5bk     1/1     Running   0          8d
```

Or better even...

```bash
# show only pods whose name starts with `frontend` or `redis`
# and only display the columns `namespace` and `name`
$ kubectl get pods --all-namespaces | tabloid \
>   --expr '(name =~ "^frontend" || name =~ "^redis") && namespace == "team-a-apps"' \
>   --column namespace,name
NAMESPACE     NAME
team-a-apps   frontend-5c6c94684f-5kzbk
team-a-apps   frontend-5c6c94684f-k2d7d
team-a-apps   frontend-5c6c94684f-ppgkx
team-a-apps   redis-follower-dddfbdcc9-9xd8l
team-a-apps   redis-follower-dddfbdcc9-l9ngl
team-a-apps   redis-leader-fb76b4755-6t5bk
```

Or we can also reorder the output:

```bash
# show only pods whose name starts with `frontend` or `redis`
# and only display the columns `namespace` and `name`, but reverse
$ kubectl get pods --all-namespaces | tabloid \
>   --expr '(name =~ "^frontend" || name =~ "^redis") && namespace == "team-a-apps"' \
>   --column name, namespace
NAME                             NAMESPACE
frontend-5c6c94684f-5kzbk        team-a-apps
frontend-5c6c94684f-k2d7d        team-a-apps
frontend-5c6c94684f-ppgkx        team-a-apps
redis-follower-dddfbdcc9-9xd8l   team-a-apps
redis-follower-dddfbdcc9-l9ngl   team-a-apps
redis-leader-fb76b4755-6t5bk     team-a-apps
```

## Features

The following features are available:

* [Column titles are always on by default](docs/column-titles.md#column-titles-always-on-by-default) and their titles are [normalized for querying with the expression language](docs/column-titles.md#column-title-normalization). Additionally, [columns can be reordered](docs/column-titles.md#column-selection-and-reordering).
* There's a [powerful expression filtering](docs/expressions.md#powerful-expression-evaluator) with [several additional built-in functions](docs/expressions.md#expression-functions) to handle specific filtering (like `kubectl` durations or pod restart count).
* Extra whitespaces (like the one that `awk` or `grep` could produce) [is removed automatically, and space count is recalculated](docs/qol-improvements.md#cleaning-up-extra-whitespace).

## Why creating this app? Isn't `enter-tool-here` enough?

The answer is "maybe". In short, I wanted to create a tool that serves my own purpose, with a quick and easy to use interface where I don't have to remember either cryptic languages or need to hack my way through to get the outputs I want.

While it's possible for `kubectl`, for example, to output JSON or YAML and have that parsed instead, I want this tool to be a one-size-fits-most in terms of column parsing. I build my own tools around the same premise of the 3+ tab padding and using Go's amazing `tabwriter`, so why not make this tool work with future versions of my own apps and potentially other 3rd-party apps?

## You have a bug, can I fix it?

Absolutely! This was a weekend project and really doesn't have much testing. Parsing columns might sound like a simple task, but you see, given the following input to the best tool out there to parse columns, `awk`, you'll see how quickly it goes wrong:

```
NAMESPACE   NAME (PROVIDED)                       READY   STATUS    RESTARTS   AGE
argocd      argocd-application-controller-0       1/1     Running   0          8d
argocd      argocd-dex-server-6dcf645b6b-nf2xb    1/1     Running   0          8d
argocd      argocd-redis-5b6967fdfc-48z9d         1/1     Running   0          8d
argocd      argocd-repo-server-7598bf5999-jfqlt   1/1     Running   0          8d
argocd      argocd-server-79f9bc9b44-5fdsp        1/1     Running   0          8d
```

```
$ awk '{ print $2 }' pods-wrong-title.txt
NAME
argocd-application-controller-0
argocd-dex-server-6dcf645b6b-nf2xb
argocd-redis-5b6967fdfc-48z9d
argocd-repo-server-7598bf5999-jfqlt
argocd-server-79f9bc9b44-5fdsp
```

The name of the 2nd column is `NAME (PROVIDED)`, yet `awk` parsed it as just `NAME`. `awk` is suitable for more generic approaches, while this tool works in harmony with `tabwriter` outputs, and as such, we can totally parse the column well:

```bash
$ cat pods-wrong-title.txt | tabloid --column name_provided
#                                 or --column "NAME (PROVIDED)"
#                                 or --column "name (provided)"
NAME (PROVIDED)
argocd-application-controller-0
argocd-dex-server-6dcf645b6b-nf2xb
argocd-redis-5b6967fdfc-48z9d
argocd-repo-server-7598bf5999-jfqlt
argocd-server-79f9bc9b44-5fdsp
```

Back to the point at hand though... Absolutely! Feel free to send any PRs you might want to see fixed/improved.
