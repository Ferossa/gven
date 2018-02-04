GVen
====

golang vendoring tool

Usage:

1. Init project - creates empty project config
   ```gven init```

2. Build targets
   ```gven build [target]```
   if no target specified - build all

3. Add dependency
   ```gven require [-t=targets] [-d=true] package[:version] [repository] [type]```
   if no targets specified - add to all
   if d=true - add to dev dependencies

4. Update dependencies
   ```gven update [packages]```
