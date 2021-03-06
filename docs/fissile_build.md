## fissile build

Has subcommands to build all images and necessary artifacts.

### Synopsis


Has subcommands to build all images and necessary artifacts.

### Options inherited from parent commands

```
  -c, --cache-dir string         Local BOSH cache directory. (default "~/.bosh/cache")
      --config string            config file (default is $HOME/.fissile.yaml)
  -f, --configgin string         Path to the tarball containing configgin.
  -d, --dark-opinions string     Path to a BOSH deployment manifest file that contains properties that should not have opinionated defaults.
  -l, --light-opinions string    Path to a BOSH deployment manifest file that contains properties to be used as defaults.
  -o, --output string            Choose output format, one of human, json, or yaml (currently only for 'show properties') (default "human")
  -r, --release string           Path to dev BOSH release(s).
  -n, --release-name string      Name of a dev BOSH release; if empty, default configured dev release name will be used
  -v, --release-version string   Version of a dev BOSH release; if empty, the latest dev release will be used
  -p, --repository string        Repository name prefix used to create image names. (default "fissile")
  -m, --role-manifest string     Path to a yaml file that details which jobs are used for each role.
  -w, --work-dir string          Path to the location of the work directory. (default "/var/fissile")
  -W, --workers int              Number of workers to use. (default 2)
```

### SEE ALSO
* [fissile](fissile.md)	 - The BOSH disintegrator
* [fissile build images](fissile_build_images.md)	 - Builds Docker images from your BOSH releases.
* [fissile build layer](fissile_build_layer.md)	 - Has subcommands for building Docker layers used during the creation of your images.
* [fissile build packages](fissile_build_packages.md)	 - Builds BOSH packages in a Docker container.

###### Auto generated by spf13/cobra on 25-Oct-2016
