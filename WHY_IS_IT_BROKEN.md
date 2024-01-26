# Common problems for why things are broken

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
