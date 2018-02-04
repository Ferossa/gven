GVen
====

golang vendoring tool

### Installing

```
go get github.com/ferossa/gven
```

### Available commands

##### init

```
gven init
```

Creates empty gven config file.

##### require

```
gven require [-t=targets] [-d=true] package[:version] [repository] [type]
```

Adds dependency to targets\
-t - comma separated targets list; if no flag specified dependency will be added to all targets\
-d - if true dependecy will be added under dev section

##### update

```
gven update [packages]
```

Update dependencies.

##### build

```
gven build [target]
```

Build targets. Build all if no targets specified.

### Notes

Development environment is determined by GODEV environment variable.