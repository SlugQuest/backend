# Common problems for why things are broken

## Most common issue: database inconsistencies

Reset the DB by:

```bash
rm slugquest.db
sqlite3 slugquest.db < schema.sql
```

## My nix fails to build

### You installed a new go dependency

If you get an error message like:
```bash
error: builder for '/nix/store/hbd4xfbmw0cfzsa1d52lglx4rksc5l62-backend-0.1.drv' failed with exit code 1;                      
       last 9 log lines:                                                                                                       
       > Running phase: unpackPhase                                                                                            
       > unpacking source archive /nix/store/1an6ihkd21is4ah01q23vd0cfqfl9qar-wss4bzj085pc6fcixsyzpcsip6vvfk1b-source          
       > source root is wss4bzj085pc6fcixsyzpcsip6vvfk1b-source                                             
       > Running phase: patchPhase                                                                                                                                                                                                                            
       > Running phase: updateAutotoolsGnuConfigScriptsPhase                                                                   
       > Running phase: configurePhase                                                                                         
       > Running phase: buildPhase                                                                                             
       > Building subPackage .                                                                                                 
       > main.go:3:8: cannot find module providing package github.com/gin-gonic/gin: import lookup disabled by -mod=vendor
       For full logs, run 'nix log /nix/store/hbd4xfbmw0cfzsa1d52lglx4rksc5l62-backend-0.1.drv'.
```

You need to enter the nix dev shell with `nix develop` and then run `gomod2nix`

### Can't get a development shell with `nix develop`

Getting something like this after doing `nix develop`:
```
error:
       … while calling the 'derivationStrict' builtin

         at /derivation-internal.nix:9:12:

            8|
            9|   strict = derivationStrict drvAttrs;
             |            ^
           10|

       … while evaluating derivation 'nix-shell'
         whose name attribute is located at /nix/store/kwd6lmx004rkv2r00vj3fcg5ijfvnagk-source/pkgs/stdenv/generic/make-derivation.nix:352:7

       … while evaluating attribute 'nativeBuildInputs' of derivation 'nix-shell'

         at /nix/store/kwd6lmx004rkv2r00vj3fcg5ijfvnagk-source/pkgs/stdenv/generic/make-derivation.nix:396:7:

          395|       depsBuildBuild              = elemAt (elemAt dependencies 0) 0;
          396|       nativeBuildInputs           = elemAt (elemAt dependencies 0) 1;
             |       ^
          397|       depsBuildTarget             = elemAt (elemAt dependencies 0) 2;

       (stack trace truncated; use '--show-trace' to show the full trace)

, but no compatible Go attribute could be found.
```

1. Run `nix flake update` -> updates `flake.lock`
